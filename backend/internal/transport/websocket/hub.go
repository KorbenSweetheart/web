// Package websocket provides WebSocket connection lifecycle management and real-time event routing.
package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
	"sync"

	"github.com/go-playground/validator/v10"
)

// ChatManager defines chat operations required by the WebSocket Hub.
type ChatManager interface {
	SaveMessage(ctx context.Context, senderID, chatID int64, content string) (*domain.Message, int64, error)
	MarkAsRead(ctx context.Context, readerID, chatID int64) error
	ChatParticipants(ctx context.Context, chatID int64) (userOneID, userTwoID int64, err error)
}

// PresenceManager defines user presence tracking required by the WebSocket Hub.
type PresenceManager interface {
	UserConnected(ctx context.Context, userID int64, connID string) (isFirst bool, err error)
	UserDisconnected(ctx context.Context, userID int64, connID string) (isLast bool, err error)
	BatchIsOnline(ctx context.Context, userIDs []int64) (map[int64]bool, error)
}

// Hub coordinates all active WebSocket client connections and event routing.
type Hub struct {
	mu          sync.RWMutex
	clients     map[int64]map[string]*Client // userID -> connID -> *Client
	Register    chan *Client
	Unregister  chan *Client
	Inbound     chan *ClientInboundMessage
	done        chan struct{}
	chatService ChatManager
	presence    PresenceManager
	validator   *validator.Validate
	log         *slog.Logger
}

// NewHub creates an instance of Hub with its dependencies.
func NewHub(chatService ChatManager, presence PresenceManager, validator *validator.Validate, logger *slog.Logger) *Hub {
	return &Hub{
		clients:     make(map[int64]map[string]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Inbound:     make(chan *ClientInboundMessage, 128),
		done:        make(chan struct{}),
		chatService: chatService,
		presence:    presence,
		validator:   validator,
		log:         logger,
	}
}

// Run executes the central Hub event loop in a dedicated goroutine.
func (h *Hub) Run(ctx context.Context) {
	defer close(h.done)

	for {
		select {
		case <-ctx.Done():
			h.cleanupAllClients()
			return

		case client := <-h.Register:
			h.registerClient(ctx, client)

		case client := <-h.Unregister:
			h.unregisterClient(ctx, client)

		case msg := <-h.Inbound:
			h.handleInboundMessage(ctx, msg)
		}
	}
}

func (h *Hub) registerClient(ctx context.Context, client *Client) {
	h.mu.Lock()
	userConns, exists := h.clients[client.UserID]
	if !exists {
		userConns = make(map[string]*Client)
		h.clients[client.UserID] = userConns
	}
	userConns[client.ConnID] = client
	h.mu.Unlock()

	h.log.Debug("client connected to websocket", "user_id", client.UserID, "conn_id", client.ConnID)

	isFirst, err := h.presence.UserConnected(ctx, client.UserID, client.ConnID)
	if err != nil {
		h.log.Debug("failed to register presence for user", "user_id", client.UserID, logger.Err(err))
		return
	}

	if isFirst {
		h.log.Debug("user transitioned to online", "user_id", client.UserID)
	}
}

func (h *Hub) unregisterClient(ctx context.Context, client *Client) {
	h.mu.Lock()
	userConns, exists := h.clients[client.UserID]
	var removed bool
	if exists {
		if _, ok := userConns[client.ConnID]; ok {
			delete(userConns, client.ConnID)
			close(client.Send)
			removed = true
			if len(userConns) == 0 {
				delete(h.clients, client.UserID)
			}
		}
	}
	h.mu.Unlock()

	if !removed {
		return
	}

	h.log.Debug("client disconnected from websocket", "user_id", client.UserID, "conn_id", client.ConnID)

	isLast, err := h.presence.UserDisconnected(ctx, client.UserID, client.ConnID)
	if err != nil {
		h.log.Debug("failed to deregister presence for user", "user_id", client.UserID, logger.Err(err))
		return
	}

	if isLast {
		h.log.Debug("user transitioned to offline", "user_id", client.UserID)
	}
}

func (h *Hub) handleInboundMessage(ctx context.Context, msg *ClientInboundMessage) {
	switch msg.Event.Type {
	case EventChatMessage:
		h.handleChatMessage(ctx, msg)
	case EventChatTyping:
		h.handleChatTyping(ctx, msg)
	case EventChatRead:
		h.handleChatRead(ctx, msg)
	case EventPresenceCheck:
		h.handlePresenceCheck(ctx, msg)
	default:
		msg.Client.SendError("unknown event type: " + msg.Event.Type)
	}
}

func (h *Hub) handleChatMessage(ctx context.Context, msg *ClientInboundMessage) {
	var payload SendMessagePayload
	if err := json.Unmarshal(msg.Event.Payload, &payload); err != nil {
		msg.Client.SendError("invalid message payload")
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		msg.Client.SendError("validation error: " + err.Error())
		return
	}

	savedMsg, recipientID, err := h.chatService.SaveMessage(ctx, msg.Client.UserID, payload.ChatID, payload.Content)
	if err != nil {
		msg.Client.SendError(err.Error())
		return
	}

	outEvent := OutboundEvent{
		Type: EventChatMessage,
		Payload: MessagePayload{
			ID:        savedMsg.ID,
			ChatID:    savedMsg.ChatID,
			SenderID:  savedMsg.SenderID,
			Content:   savedMsg.Content,
			CreatedAt: savedMsg.CreatedAt,
			IsViewed:  savedMsg.IsViewed,
		},
	}

	// 1. Deliver ACK/saved event back to sender (all active tabs)
	h.deliverToUser(msg.Client.UserID, outEvent)

	// 2. Deliver event to recipient if currently online
	if recipientID > 0 {
		h.deliverToUser(recipientID, outEvent)
	}
}

func (h *Hub) handleChatTyping(ctx context.Context, msg *ClientInboundMessage) {
	var payload TypingPayload
	if err := json.Unmarshal(msg.Event.Payload, &payload); err != nil {
		msg.Client.SendError("invalid typing payload")
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		msg.Client.SendError("validation error: " + err.Error())
		return
	}

	recipientID := h.resolveOtherParticipant(ctx, payload.ChatID, msg.Client.UserID)
	if recipientID > 0 {
		h.deliverToUser(recipientID, OutboundEvent{
			Type: EventChatTyping,
			Payload: TypingBroadcastPayload{
				ChatID:   payload.ChatID,
				UserID:   msg.Client.UserID,
				IsTyping: payload.IsTyping,
			},
		})
	}
}

func (h *Hub) handleChatRead(ctx context.Context, msg *ClientInboundMessage) {
	var payload ReadMessagesPayload
	if err := json.Unmarshal(msg.Event.Payload, &payload); err != nil {
		msg.Client.SendError("invalid read payload")
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		msg.Client.SendError("validation error: " + err.Error())
		return
	}

	if err := h.chatService.MarkAsRead(ctx, msg.Client.UserID, payload.ChatID); err != nil {
		msg.Client.SendError(err.Error())
		return
	}

	recipientID := h.resolveOtherParticipant(ctx, payload.ChatID, msg.Client.UserID)
	if recipientID > 0 {
		h.deliverToUser(recipientID, OutboundEvent{
			Type: EventChatRead,
			Payload: ReadBroadcastPayload{
				ChatID:   payload.ChatID,
				ReaderID: msg.Client.UserID,
			},
		})
	}
}

func (h *Hub) handlePresenceCheck(ctx context.Context, msg *ClientInboundMessage) {
	var payload PresenceCheckPayload
	if err := json.Unmarshal(msg.Event.Payload, &payload); err != nil {
		msg.Client.SendError("invalid presence check payload")
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		msg.Client.SendError("validation error: " + err.Error())
		return
	}

	statuses, err := h.presence.BatchIsOnline(ctx, payload.UserIDs)
	if err != nil {
		msg.Client.SendError(err.Error())
		return
	}

	msg.Client.SendEvent(OutboundEvent{
		Type: EventPresenceBatch,
		Payload: PresenceBatchPayload{
			Statuses: statuses,
		},
	})
}

func (h *Hub) resolveOtherParticipant(ctx context.Context, chatID, senderID int64) int64 {
	userOneID, userTwoID, err := h.chatService.ChatParticipants(ctx, chatID)
	if err != nil {
		h.log.Debug("failed to resolve chat participants", "chat_id", chatID, logger.Err(err))
		return 0
	}

	if senderID == userOneID {
		return userTwoID
	}
	return userOneID
}

func (h *Hub) deliverToUser(userID int64, event OutboundEvent) {
	h.mu.RLock()
	userConns, exists := h.clients[userID]
	if !exists || len(userConns) == 0 {
		h.mu.RUnlock()
		return
	}

	targetClients := make([]*Client, 0, len(userConns))
	for _, client := range userConns {
		targetClients = append(targetClients, client)
	}
	h.mu.RUnlock()

	for _, client := range targetClients {
		client.SendEvent(event)
	}
}

func (h *Hub) cleanupAllClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, userConns := range h.clients {
		for _, client := range userConns {
			close(client.Send)
		}
		delete(h.clients, userID)
	}
}
