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
	FindDirectChat(ctx context.Context, userA, userB int64) (*domain.Chat, error)
	CreateDirectChat(ctx context.Context, chat *domain.Chat) error
	FindChatByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	FindUserChats(ctx context.Context, userID int64) ([]*domain.Chat, error)
	SaveMessage(ctx context.Context, message *domain.Message) error
	LoadChatHistory(ctx context.Context, chatID, lastMessageID int64, limit int) ([]*domain.Message, error)
	MarkMessagesAsRead(ctx context.Context, chatID, readerID int64) error
}

// UserConnectionChecker verifies relationship status between users.
type UserConnectionChecker interface {
	FindConnectionRecord(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error)
}

// ChatService coordinates chat business rules and interactions between users.
type ChatService struct {
	chatRepo ChatRepository
	connRepo UserConnectionChecker
	log      *slog.Logger
}

// NewChatService creates an instance of ChatService with repository dependencies.
func NewChatService(cr ChatRepository, cc UserConnectionChecker, logger *slog.Logger) *ChatService {
	return &ChatService{
		chatRepo: cr,
		connRepo: cc,
		log:      logger,
	}
}

// DirectChat retrieves an existing chat between two users or creates a new one
// provided both users have an accepted connection.
func (s *ChatService) DirectChat(ctx context.Context, requesterID, targetUserID int64) (*domain.Chat, error) {
	const op = "service.chatService.DirectChat"

	if requesterID == targetUserID {
		return nil, fmt.Errorf("%s: cannot create chat with self: %d", op, requesterID)
	}

	conn, err := s.connRepo.FindConnectionRecord(ctx, requesterID, targetUserID)
	if err != nil {
		if errors.Is(err, domain.ErrConnectionNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrUsersNotConnected)
		}
		return nil, fmt.Errorf("%s: failed to verify connection between %d and %d: %w", op, requesterID, targetUserID, err)
	}

	if conn.Status != domain.Accepted {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUsersNotConnected)
	}

	userOneID, userTwoID := domain.NormalizeUserPair(requesterID, targetUserID)

	chat, err := s.chatRepo.FindDirectChat(ctx, userOneID, userTwoID)
	if err == nil {
		return chat, nil
	}

	if !errors.Is(err, domain.ErrChatNotFound) {
		return nil, fmt.Errorf("%s: failed to find direct chat: %w", op, err)
	}

	newChat := &domain.Chat{
		UserOneID: userOneID,
		UserTwoID: userTwoID,
	}

	if err := s.chatRepo.CreateDirectChat(ctx, newChat); err != nil {
		return nil, fmt.Errorf("%s: failed to create direct chat: %w", op, err)
	}

	return newChat, nil
}

// ChatHistory loads a paginated batch of messages for a chat after verifying user membership.
func (s *ChatService) ChatHistory(ctx context.Context, requesterID, chatID, lastMessageID int64, limit int) ([]*domain.Message, error) {
	const op = "service.chatService.ChatHistory"

	if err := s.verifyChatMembership(ctx, requesterID, chatID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if limit <= 0 || limit > MaxMessageHistoryLimit {
		limit = DefaultMessageHistoryLimit
	}

	messages, err := s.chatRepo.LoadChatHistory(ctx, chatID, lastMessageID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

// UserChats returns all active conversations for the user.
func (s *ChatService) UserChats(ctx context.Context, requesterID int64) ([]*domain.Chat, error) {
	const op = "service.chatService.UserChats"

	chats, err := s.chatRepo.FindUserChats(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user chats: %w", op, err)
	}

	return chats, nil
}

// SaveMessage validates message constraints, checks sender membership, and persists the message.
func (s *ChatService) SaveMessage(ctx context.Context, senderID, chatID int64, content string) (*domain.Message, error) {
	const op = "service.chatService.SaveMessage"

	if err := s.verifyChatMembership(ctx, senderID, chatID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("%s: message content cannot be empty", op)
	}

	if len(content) > MaxMessageLength {
		return nil, fmt.Errorf("%s: message content exceeds maximum length of %d characters", op, MaxMessageLength)
	}

	msg := &domain.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}

	if err := s.chatRepo.SaveMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return msg, nil
}

// MarkAsRead marks unread incoming messages in a chat as viewed by readerID.
func (s *ChatService) MarkAsRead(ctx context.Context, readerID, chatID int64) error {
	const op = "service.chatService.MarkAsRead"

	if err := s.verifyChatMembership(ctx, readerID, chatID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.chatRepo.MarkMessagesAsRead(ctx, chatID, readerID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ChatParticipants returns the two participant user IDs for a given chat.
func (s *ChatService) ChatParticipants(ctx context.Context, chatID int64) (userOneID, userTwoID int64, err error) {
	const op = "service.chatService.ChatParticipants"

	chat, err := s.chatRepo.FindChatByID(ctx, chatID)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return chat.UserOneID, chat.UserTwoID, nil
}

// verifyChatMembership checks whether a user belongs to the specified chat.
func (s *ChatService) verifyChatMembership(ctx context.Context, userID, chatID int64) error {
	chat, err := s.chatRepo.FindChatByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.UserOneID != userID && chat.UserTwoID != userID {
		return domain.ErrNotChatParticipant
	}

	return nil
}
