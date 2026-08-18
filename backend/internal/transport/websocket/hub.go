package websocket

import (
	"context"
	"log/slog"
	"match-me-api/internal/logger"
)

// https://echo.labstack.com/cookbook/websocket/

type Chat struct {
	ID      int64             `json:"id"`
	Clients map[int64]*Client `json:"clients"`
	History []*Message
}

type Hub struct {
	Clients    map[*Client]bool
	Chats      map[int64]*Chat // map[int64]*Chat or map[int64]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	log        *slog.Logger
}

func NewCore(logger *slog.Logger) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Chats:      make(map[int64]*Chat),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		log:        logger,
	}
}

// The core will be ran in a different go Routine
func (hub *Hub) Run(ctx context.Context) {
	const op = "service.matchService.Recommendations"
	log := hub.log.With(slog.String("op", op)) // TODO: Check that it is a proper way to use it here, no race condition

	for {
		select {
		case cl := <-hub.Register:
			// add client to hub list
			if chat, ok := hub.Chats[cl.ChatID]; ok {
				if _, ok := chat.Clients[cl.ID]; !ok {
					chat.Clients[cl.ID] = cl
				}
				go func() {
					messages, err := hub.chatRepo.LoadChatHistory(ctx, chat.ID, 10)
					if err != nil {
						log.Debug("failed to load chat history", "id", chat.ID, "error", logger.Err(err))
						return
					}

					for _, msg := range messages {
						wsMsg := &Message{
							Content:   msg.Content,
							ChatID:    cl.ChatID,
							Username:  msg.Sender.Name,
							SenderID:  msg.SenderID,
							Timestamp: msg.CreatedAt,
						}
						cl.Message <- wsMsg
					}
				}()
			}

		case cl := <-hub.Unregister:
			if _, ok := hub.Chats[cl.ChatID]; ok {
				if _, ok := hub.Chats[cl.ChatID].Clients[cl.ID]; ok {
					delete(hub.Chats[cl.ChatID].Clients, cl.ID)
					close(cl.Message)
				}
			}

			// FAN OUT
		case m := <-hub.Broadcast:
			if room, ok := hub.Chats[m.ChatID]; ok {
				room.History = append(room.History, m)

				go func(msg *Message) {

					var SenderID int64

					dbMsg := &roomRepo.Message{
						RoomID:   roomUUID,
						UserID:   userID,
						Username: msg.Username,
						Content:  msg.Content,
						IsSystem: msg.System,
					}

					if _, err := hub.roomRepo.CreateMessage(context.Background(), dbMsg); err != nil {
						// log.Printf("Failed to persist message: %v", err)
					}

					if userID != nil {
						if err := hub.statsRepo.IncrementMessageCount(context.Background(), *userID); err != nil {
							// log.Printf("Failed to update message count for user %s: %v", userID.String(), err)
						}
					}
				}(m)

				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
