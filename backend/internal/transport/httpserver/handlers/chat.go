package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// ChatManager defines business methods required by ChatHandler.
type ChatManager interface {
	DirectChat(ctx context.Context, requesterID, targetUserID int64) (*domain.Chat, error)
	ChatHistory(ctx context.Context, requesterID, chatID, lastMessageID int64, limit int) ([]*domain.Message, error)
	UserChats(ctx context.Context, requesterID int64) ([]*domain.Chat, error)
}

// ChatHandler handles HTTP REST endpoints for chats and historical messages.
type ChatHandler struct {
	chatService ChatManager
	validator   *validator.Validate
	log         *slog.Logger
}

// NewChatHandler constructs a new ChatHandler instance.
func NewChatHandler(cm ChatManager, v *validator.Validate, l *slog.Logger) *ChatHandler {
	return &ChatHandler{
		chatService: cm,
		validator:   v,
		log:         l,
	}
}

// DirectChat retrieves or creates a 1-on-1 direct chat with another user.
// POST /chats/direct
func (h *ChatHandler) DirectChat(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	var req dto.DirectChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	chat, err := h.chatService.DirectChat(ctx, myID, req.TargetUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUsersNotConnected) {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": domain.ErrUsersNotConnected.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to open direct chat"})
	}

	return c.JSON(http.StatusOK, formatChatResponse(chat))
}

// UserChats returns all direct chats belonging to the authenticated user.
// GET /chats
func (h *ChatHandler) UserChats(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	chats, err := h.chatService.UserChats(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to load user chats"})
	}

	response := make([]dto.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		response = append(response, formatChatResponse(chat))
	}

	return c.JSON(http.StatusOK, response)
}

// ChatHistory returns paginated historical messages for a given chat.
// GET /chats/:id/messages?limit=15&last_message_id=100
func (h *ChatHandler) ChatHistory(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil || chatID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid chat id: NAN"})
	}

	limit := 0
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	var lastMessageID int64
	if lastMsgStr := c.QueryParam("last_message_id"); lastMsgStr != "" {
		if parsedLastMsgID, err := strconv.ParseInt(lastMsgStr, 10, 64); err == nil {
			lastMessageID = parsedLastMsgID
		}
	}

	messages, err := h.chatService.ChatHistory(ctx, myID, chatID, lastMessageID, limit)
	if err != nil {
		if errors.Is(err, domain.ErrChatNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{"error": domain.ErrChatNotFound.Error()})
		}
		if errors.Is(err, domain.ErrNotChatParticipant) {
			return c.JSON(http.StatusForbidden, map[string]any{"error": domain.ErrNotChatParticipant.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to load chat history"})
	}

	response := make([]dto.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		response = append(response, dto.MessageResponse{
			ID:        msg.ID,
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			IsViewed:  msg.IsViewed,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func formatChatResponse(chat *domain.Chat) dto.ChatResponse {
	return dto.ChatResponse{
		ID:        chat.ID,
		UserOneID: chat.UserOneID,
		UserTwoID: chat.UserTwoID,
		CreatedAt: chat.CreatedAt,
		UserOne: dto.UserSummaryResponse{
			ID:         chat.UserOne.UserID,
			Name:       chat.UserOne.Name,
			PictureURL: chat.UserOne.PictureURL,
		},
		UserTwo: dto.UserSummaryResponse{
			ID:         chat.UserTwo.UserID,
			Name:       chat.UserTwo.Name,
			PictureURL: chat.UserTwo.PictureURL,
		},
	}
}
