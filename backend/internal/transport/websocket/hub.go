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
	SaveMessage(ctx context.Context, senderID, chatID int64, content string) (*domain.Message, error)
	MarkAsRead(ctx context.Context, readerID, chatID int64) error
	ChatParticipants(ctx context.Context, chatID int64) (userOneID, userTwoID int64, err error)
}

// ConnectionProvider defines connection queries required for presence broadcasts.
type ConnectionProvider interface {
	AcceptedConnections(ctx context.Context, userID int64) ([]int64, error)
}

// Hub coordinates all active WebSocket client connections, presence state, and event routing.
type Hub struct {
	mu          sync.RWMutex
	clients     map[int64]*Client
	Register    chan *Client
	Unregister  chan *Client
	Inbound     chan *ClientInboundMessage
	chatService ChatManager
	connService ConnectionProvider
	validator   *validator.Validate
	log         *slog.Logger
}

// NewHub creates an instance of Hub with its service dependencies.
func NewHub(chatService ChatManager, connService ConnectionProvider, validator *validator.Validate, logger *slog.Logger) *Hub {
	return &Hub{
		clients:     make(map[int64]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Inbound:     make(chan *ClientInboundMessage, 128),
		chatService: chatService,
		connService: connService,
		validator:   validator,
		log:         logger,
	}
}

// Run executes the central Hub event loop in a dedicated goroutine.
func (h *Hub) Run(ctx context.Context) {
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

// IsOnline checks whether a given user is currently connected to this Hub.
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, online := h.clients[userID]
	return online
}

func (h *Hub) registerClient(ctx context.Context, client *Client) {
	h.mu.Lock()
	h.clients[client.UserID] = client
	h.mu.Unlock()

	h.log.Debug("client connected to websocket", "user_id", client.UserID)

	// Broadcast online presence to all accepted connections
	h.broadcastUserStatus(ctx, client.UserID, "online")
}

func (h *Hub) unregisterClient(ctx context.Context, client *Client) {
	h.mu.Lock()
	existing, ok := h.clients[client.UserID]
	if ok && existing == client {
		delete(h.clients, client.UserID)
		close(client.Send)
	}
	h.mu.Unlock()

	if ok && existing == client {
		h.log.Debug("client disconnected from websocket", "user_id", client.UserID)
		// Broadcast offline presence to all accepted connections
		h.broadcastUserStatus(ctx, client.UserID, "offline")
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

	savedMsg, err := h.chatService.SaveMessage(ctx, msg.Client.UserID, payload.ChatID, payload.Content)
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

	// 1. Deliver ACK/saved event back to sender
	msg.Client.SendEvent(outEvent)

	// 2. Deliver event to recipient if currently online
	recipientID := h.resolveOtherParticipant(ctx, payload.ChatID, msg.Client.UserID)
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
	client, online := h.clients[userID]
	h.mu.RUnlock()

	if online {
		client.SendEvent(event)
	}
}

func (h *Hub) broadcastUserStatus(ctx context.Context, userID int64, status string) {
	connectedUserIDs, err := h.connService.AcceptedConnections(ctx, userID)
	if err != nil {
		h.log.Debug("failed to fetch accepted connections for presence update", "user_id", userID, logger.Err(err))
		return
	}

	event := OutboundEvent{
		Type: EventUserStatus,
		Payload: UserStatusPayload{
			UserID: userID,
			Status: status,
		},
	}

	for _, friendID := range connectedUserIDs {
		h.deliverToUser(friendID, event)
	}
}

func (h *Hub) cleanupAllClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, client := range h.clients {
		delete(h.clients, userID)
		close(client.Send)
	}
}
