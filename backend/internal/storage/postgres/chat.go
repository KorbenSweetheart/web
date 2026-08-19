package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm"
)

// CreateDirectChat inserts a new 1-on-1 direct chat record into the database.
func (s *Storage) CreateDirectChat(ctx context.Context, chat *domain.Chat) error {
	const op = "storage.postgres.CreateDirectChat"

	chat.UserOneID, chat.UserTwoID = domain.NormalizeUserPair(chat.UserOneID, chat.UserTwoID)

	if err := s.db.WithContext(ctx).Create(chat).Error; err != nil {
		return fmt.Errorf("%s: failed to create direct chat: %w", op, err)
	}

	return nil
}

// FindDirectChat finds an existing 1-on-1 direct chat record between two users.
func (s *Storage) FindDirectChat(ctx context.Context, userA, userB int64) (*domain.Chat, error) {
	const op = "storage.postgres.FindDirectChat"

	userOneID, userTwoID := domain.NormalizeUserPair(userA, userB)

	var chat domain.Chat
	err := s.db.WithContext(ctx).
		Preload("UserOne").
		Preload("UserTwo").
		Where("user_one_id = ? AND user_two_id = ?", userOneID, userTwoID).
		First(&chat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrChatNotFound
		}
		return nil, fmt.Errorf("%s: failed to find direct chat between %d and %d: %w", op, userA, userB, err)
	}

	return &chat, nil
}

// FindChatByID retrieves a single chat by its primary ID.
func (s *Storage) FindChatByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	const op = "storage.postgres.FindChatByID"

	var chat domain.Chat
	err := s.db.WithContext(ctx).
		Preload("UserOne").
		Preload("UserTwo").
		Where("id = ?", chatID).
		First(&chat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrChatNotFound
		}
		return nil, fmt.Errorf("%s: failed to find chat by id %d: %w", op, chatID, err)
	}

	return &chat, nil
}

// FindUserChats retrieves all active chats for a user, ordered by latest activity.
func (s *Storage) FindUserChats(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	const op = "storage.postgres.FindUserChats"

	chats := make([]*domain.Chat, 0)
	err := s.db.WithContext(ctx).
		Preload("UserOne").
		Preload("UserTwo").
		Where("user_one_id = ? OR user_two_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&chats).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to find user chats for user %d: %w", op, userID, err)
	}

	return chats, nil
}

// SaveMessage persists a new chat message to the database.
func (s *Storage) SaveMessage(ctx context.Context, message *domain.Message) error {
	const op = "storage.postgres.SaveMessage"

	if err := s.db.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("%s: failed to save message: %w", op, err)
	}

	return nil
}

// LoadChatHistory loads a paginated batch of messages for a chat, descending from lastMessageID cursor.
func (s *Storage) LoadChatHistory(ctx context.Context, chatID, lastMessageID int64, limit int) ([]*domain.Message, error) {
	const op = "storage.postgres.LoadChatHistory"

	query := s.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Preload("Sender")

	if lastMessageID > 0 {
		query = query.Where("id < ?", lastMessageID)
	}

	messages := make([]*domain.Message, 0)
	err := query.
		Order("id DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to load chat history for chat %d: %w", op, chatID, err)
	}

	return messages, nil
}

// MarkMessagesAsRead marks unread incoming messages in a chat as viewed.
func (s *Storage) MarkMessagesAsRead(ctx context.Context, chatID, readerID int64) error {
	const op = "storage.postgres.MarkMessagesAsRead"

	err := s.db.WithContext(ctx).
		Model(&domain.Message{}).
		Where("chat_id = ? AND sender_id != ? AND is_viewed = ?", chatID, readerID, false).
		Update("is_viewed", true).Error

	if err != nil {
		return fmt.Errorf("%s: failed to mark messages as read in chat %d: %w", op, chatID, err)
	}

	return nil
}
