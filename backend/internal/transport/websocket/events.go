package websocket

import (
	"encoding/json"
	"time"
)

// Event type constants for the WebSocket wire protocol.
const (
	EventChatMessage   = "chat:message"
	EventChatTyping    = "chat:typing"
	EventChatRead      = "chat:read"
	EventPresenceCheck = "presence:check"
	EventPresenceBatch = "presence:batch"
	EventError         = "error"
)

// InboundEvent is the generic envelope for all messages sent from client to server.
// Using json.RawMessage allows delayed deserialization of the Payload based on Type.
type InboundEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// OutboundEvent is the generic envelope for all messages sent from server to client.
type OutboundEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// SendMessagePayload is received when a user sends a chat message.
type SendMessagePayload struct {
	ChatID  int64  `json:"chat_id" validate:"required,gt=0"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

// TypingPayload is received when a user starts or stops typing in a chat.
type TypingPayload struct {
	ChatID   int64 `json:"chat_id" validate:"required,gt=0"`
	IsTyping bool  `json:"is_typing"`
}

// ReadMessagesPayload is received when a user views a chat.
type ReadMessagesPayload struct {
	ChatID int64 `json:"chat_id" validate:"required,gt=0"`
}

// MessagePayload is sent to clients when a new message is created.
type MessagePayload struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsViewed  bool      `json:"is_viewed"`
}

// TypingBroadcastPayload notifies the recipient about the other user's typing state.
type TypingBroadcastPayload struct {
	ChatID   int64 `json:"chat_id"`
	UserID   int64 `json:"user_id"`
	IsTyping bool  `json:"is_typing"`
}

// ReadBroadcastPayload notifies the sender that their messages were read.
type ReadBroadcastPayload struct {
	ChatID   int64 `json:"chat_id"`
	ReaderID int64 `json:"reader_id"`
}

// PresenceCheckPayload is received when a client queries online status for a specific list of users.
type PresenceCheckPayload struct {
	UserIDs []int64 `json:"user_ids" validate:"required,min=1,max=100"`
}

// PresenceBatchPayload returns the presence statuses for requested user IDs.
type PresenceBatchPayload struct {
	Statuses map[int64]bool `json:"statuses"`
}

// ErrorPayload notifies the client about validation or permission errors.
type ErrorPayload struct {
	Message string `json:"message"`
}
