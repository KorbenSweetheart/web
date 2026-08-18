package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"time"

	"github.com/go-playground/validator/v10"
)

type ChatHub interface {
	SendMessage(ctx context.Context, senderID, receiverID, content string, timestamp time.Time) (*domain.Message, error)
}

type ChatHandler struct {
	chatService ChatHub
	validator   *validator.Validate
	log         *slog.Logger
}

func NewChatHandler(ch ChatHub, v *validator.Validate, l *slog.Logger) *ChatHandler {
	return &ChatHandler{chatService: ch, validator: v, log: l}
}

// ??? Do we need Upgrader?
// Echo GET /ws -> Upgrades to WS -> Spawns readPump & writePump -> Registers with Hub

// Chat Creation / Finding (POST /chats/direct)

// Load/OpenChat()
// Opening the chat fetches historical messages
// GET /chats/:id/messages?limit=50

// Sending a Live Message
// React sends JSON over the ALREADY OPEN WebSocket:
//  { "type": "chat:message", "chat_id": 42, "content": "Hey!" }

// Receiving a Live Message / Typing Indicator

// Global Unread Counters
