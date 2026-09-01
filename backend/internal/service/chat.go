package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
	"strings"
)

const (
	// DefaultMessageHistoryLimit is the default page size for message history.
	DefaultMessageHistoryLimit = 15
	// MaxMessageHistoryLimit is the upper bound for a single message history query.
	MaxMessageHistoryLimit = 50
	// MaxMessageLength is the maximum allowed characters in a chat message.
	MaxMessageLength = 2000
)

// ChatRepository defines the required database storage operations for chats.
type ChatRepository interface {
	GetOrCreateDirectChat(ctx context.Context, userOneID, userTwoID int64) (*domain.Chat, error)
	FindChatByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	FindUserChats(ctx context.Context, userID int64) ([]*domain.Chat, error)
	SaveMessage(ctx context.Context, message *domain.Message) (recipientID int64, err error)
	LoadChatHistory(ctx context.Context, chatID, userID, lastMessageID int64, limit int) ([]*domain.Message, error)
	MarkMessagesAsRead(ctx context.Context, chatID, readerID int64) error
}

// ChatService coordinates chat business rules and interactions between users.
type ChatService struct {
	repo ChatRepository
	log  *slog.Logger
}

// NewChatService creates an instance of ChatService with repository dependencies.
func NewChatService(r ChatRepository, logger *slog.Logger) *ChatService {
	return &ChatService{
		repo: r,
		log:  logger,
	}
}

// DirectChat retrieves an existing chat between two users or creates a new one
// provided both users have an accepted connection.
func (s *ChatService) DirectChat(ctx context.Context, requesterID, targetUserID int64) (*domain.Chat, error) {
	const op = "service.chatService.DirectChat"

	if requesterID == targetUserID {
		return nil, fmt.Errorf("%s: cannot create chat with self: %d", op, requesterID)
	}

	userOneID, userTwoID := domain.NormalizeUserPair(requesterID, targetUserID)

	chat, err := s.repo.GetOrCreateDirectChat(ctx, userOneID, userTwoID)
	if err != nil {
		if errors.Is(err, domain.ErrUsersNotConnected) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrUsersNotConnected)
		}
		return nil, fmt.Errorf("%s: failed to get or create direct chat: %w", op, err)
	}

	return chat, nil
}

// ChatHistory loads a paginated batch of messages for a chat scoped to user membership.
func (s *ChatService) ChatHistory(ctx context.Context, requesterID, chatID, lastMessageID int64, limit int) ([]*domain.Message, error) {
	const op = "service.chatService.ChatHistory"

	if limit <= 0 || limit > MaxMessageHistoryLimit {
		limit = DefaultMessageHistoryLimit
	}

	messages, err := s.repo.LoadChatHistory(ctx, chatID, requesterID, lastMessageID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

// UserChats returns all active conversations for the user.
func (s *ChatService) UserChats(ctx context.Context, requesterID int64) ([]*domain.Chat, error) {
	const op = "service.chatService.UserChats"

	chats, err := s.repo.FindUserChats(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user chats: %w", op, err)
	}

	return chats, nil
}

// SaveMessage validates message constraints, checks sender membership atomically, persists the message, and returns the recipient's user ID.
func (s *ChatService) SaveMessage(ctx context.Context, senderID, chatID int64, content string) (*domain.Message, int64, error) {
	const op = "service.chatService.SaveMessage"

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, 0, fmt.Errorf("%s: message content cannot be empty", op)
	}

	if len(content) > MaxMessageLength {
		return nil, 0, fmt.Errorf("%s: message content exceeds maximum length of %d characters", op, MaxMessageLength)
	}

	msg := &domain.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}

	recipientID, err := s.repo.SaveMessage(ctx, msg)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return msg, recipientID, nil
}

// MarkAsRead marks unread incoming messages in a chat as viewed by readerID.
func (s *ChatService) MarkAsRead(ctx context.Context, readerID, chatID int64) error {
	const op = "service.chatService.MarkAsRead"

	if err := s.repo.MarkMessagesAsRead(ctx, chatID, readerID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ChatParticipants returns the two participant user IDs for a given chat.
func (s *ChatService) ChatParticipants(ctx context.Context, chatID int64) (userOneID, userTwoID int64, err error) {
	const op = "service.chatService.ChatParticipants"

	chat, err := s.repo.FindChatByID(ctx, chatID)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return chat.UserOneID, chat.UserTwoID, nil
}
