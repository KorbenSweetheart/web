package service

import (
	"context"
	"match-me-api/internal/domain"
)

// /connections
// `SELECT to_user_id FROM connections WHERE from_user_id = $1 AND status = 'accepted'`;

// Problem: on DB level no restriction on duplicate chat creation between UserA & UserB.
// Solution: add composite UNIQUE(user_one_id, user_two_id). also to prevent duplicates (A, B) and (B, A),

// no Composite Index for message sort
// Problem: field chat_id indexed, but loading messages requires pagination and message sorting based on the timestamp:
// WHERE chat_id = X ORDER BY created_at DESC. Single-column index on chat_id will force DB make additional sort (Sort Node in EXPLAIN ANALYZE).
// Solution: Add composite index gorm:"index:idx_chat_messages,priority:1" for chat_id and gorm:"index:idx_chat_messages,priority:2" for id or created_at.

type MessageRepository interface {
	SaveMessage(ctx context.Context, chatID, senderID, receiverID int64, content string) (*domain.Message, error)
	LoadChatHistory(ctx context.Context, chatID int64, limit int) ([]*domain.Message, error)
}
