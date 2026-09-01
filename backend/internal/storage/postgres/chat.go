package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

// getOrCreateDirectChatDTO is an internal storage projection DTO for scanning raw get or create direct chat query results.
type getOrCreateDirectChatDTO struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// GetOrCreateDirectChat atomically verifies connection status and inserts or retrieves a 1-on-1 direct chat record, preloading participant profiles.
func (s *Storage) GetOrCreateDirectChat(ctx context.Context, userOneID, userTwoID int64) (*domain.Chat, error) {
	const op = "storage.postgres.GetOrCreateDirectChat"

	const sqlQuery = `
INSERT INTO chats (user_one_id, user_two_id, created_at)
SELECT ?, ?, NOW()
WHERE EXISTS (
    SELECT 1 FROM connections
    WHERE ((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))
      AND status = ?
)
ON CONFLICT (user_one_id, user_two_id)
DO UPDATE SET user_one_id = EXCLUDED.user_one_id
RETURNING id, created_at;`

	var res getOrCreateDirectChatDTO

	err := s.db.WithContext(ctx).Raw(
		sqlQuery,
		userOneID, userTwoID,
		userOneID, userTwoID, userTwoID, userOneID,
		domain.Accepted,
	).Scan(&res).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute get or create direct chat: %w", op, err)
	}

	if res.ID == 0 {
		return nil, domain.ErrUsersNotConnected
	}

	var chat domain.Chat
	if err := s.db.WithContext(ctx).
		Preload("UserOne").
		Preload("UserTwo").
		First(&chat, res.ID).Error; err != nil {
		return nil, fmt.Errorf("%s: failed to preload chat profiles: %w", op, err)
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
		Where("EXISTS (SELECT 1 FROM connections WHERE ((from_user_id = chats.user_one_id AND to_user_id = chats.user_two_id) OR (from_user_id = chats.user_two_id AND to_user_id = chats.user_one_id)) AND status = ?)", domain.Accepted).
		Order("COALESCE((SELECT MAX(created_at) FROM messages WHERE messages.chat_id = chats.id), chats.created_at) DESC").
		Find(&chats).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to find user chats for user %d: %w", op, userID, err)
	}

	return chats, nil
}

// saveMessageDTO is an internal storage projection DTO for scanning raw CTE message insert results.
type saveMessageDTO struct {
	ID          int64     `gorm:"column:id"`
	ChatID      int64     `gorm:"column:chat_id"`
	SenderID    int64     `gorm:"column:sender_id"`
	Content     string    `gorm:"column:content"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	IsViewed    bool      `gorm:"column:is_viewed"`
	RecipientID int64     `gorm:"column:recipient_id"`
}

// SaveMessage persists a new chat message to the database and returns the recipient's user ID.
func (s *Storage) SaveMessage(ctx context.Context, message *domain.Message) (int64, error) {
	const op = "storage.postgres.SaveMessage"

	const sqlQuery = `
WITH target_chat AS (
    SELECT id,
           CASE WHEN user_one_id = ? THEN user_two_id ELSE user_one_id END AS recipient_id
    FROM chats
    WHERE id = ? AND (user_one_id = ? OR user_two_id = ?)
),
inserted AS (
    INSERT INTO messages (chat_id, sender_id, content, created_at, is_viewed)
    SELECT id, ?, ?, NOW(), false
    FROM target_chat
    RETURNING id, chat_id, sender_id, content, created_at, is_viewed
)
SELECT inserted.id, inserted.chat_id, inserted.sender_id, inserted.content, inserted.created_at, inserted.is_viewed, target_chat.recipient_id
FROM inserted
CROSS JOIN target_chat;`

	var res saveMessageDTO
	err := s.db.WithContext(ctx).Raw(
		sqlQuery,
		message.SenderID,
		message.ChatID,
		message.SenderID,
		message.SenderID,
		message.SenderID,
		message.Content,
	).Scan(&res).Error

	if err != nil {
		return 0, fmt.Errorf("%s: failed to save message: %w", op, err)
	}

	if res.ID == 0 {
		return 0, domain.ErrNotChatParticipant
	}

	message.ID = res.ID
	message.CreatedAt = res.CreatedAt
	message.IsViewed = res.IsViewed

	return res.RecipientID, nil
}

// LoadChatHistory loads a paginated batch of messages for a chat scoped to user membership, descending from lastMessageID cursor.
func (s *Storage) LoadChatHistory(ctx context.Context, chatID, userID, lastMessageID int64, limit int) ([]*domain.Message, error) {
	const op = "storage.postgres.LoadChatHistory"

	query := s.db.WithContext(ctx).
		Model(&domain.Message{}).
		Joins("JOIN chats ON chats.id = messages.chat_id").
		Where("messages.chat_id = ? AND (chats.user_one_id = ? OR chats.user_two_id = ?)", chatID, userID, userID).
		Preload("Sender")

	if lastMessageID > 0 {
		query = query.Where("messages.id < ?", lastMessageID)
	}

	messages := make([]*domain.Message, 0)
	err := query.
		Order("messages.id DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to load chat history for chat %d: %w", op, chatID, err)
	}

	return messages, nil
}

// MarkMessagesAsRead marks unread incoming messages in a chat as viewed by readerID.
func (s *Storage) MarkMessagesAsRead(ctx context.Context, chatID, readerID int64) error {
	const op = "storage.postgres.MarkMessagesAsRead"

	err := s.db.WithContext(ctx).
		Model(&domain.Message{}).
		Where("chat_id = ? AND sender_id != ? AND is_viewed = ?", chatID, readerID, false).
		Where("EXISTS (SELECT 1 FROM chats WHERE chats.id = messages.chat_id AND (chats.user_one_id = ? OR chats.user_two_id = ?))", readerID, readerID).
		Update("is_viewed", true).Error

	if err != nil {
		return fmt.Errorf("%s: failed to mark messages as read in chat %d: %w", op, chatID, err)
	}

	return nil
}
