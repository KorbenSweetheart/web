package dto

import (
	"time"
)

// REQUEST DTOs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ProfileUpdateRequest struct {
	Name            *string            `json:"name" validate:"omitempty,min=2,max=100"`
	PictureURL      *string            `json:"picture_url" validate:"omitempty,url"`
	Age             *int64             `json:"age" validate:"omitempty,number,gte=0,lte=125"`
	Bio             *string            `json:"bio" validate:"omitempty,max=1000"`
	MaxRadius       *float64           `json:"max_radius" validate:"omitempty,number,gte=1"`
	InteractionMode *int               `json:"interaction_mode" validate:"omitempty,number,gte=1,lte=4"`
	Activities      *[]ActivityRequest `json:"activities" validate:"omitempty,dive"`
	Lat             *float64           `json:"lat" validate:"omitempty,latitude"`
	Lon             *float64           `json:"lon" validate:"omitempty,longitude"`
}

type ActivityRequest struct {
	ID            int64 `json:"id" validate:"required,gte=1"`
	Experience    int   `json:"experience_level" validate:"required,gte=1,lte=5"` // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int   `json:"interest_level" validate:"required,gte=1,lte=5"`   // 1-5 levels: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

// RESPONSE DTOs

// /users/{id}
type UserSummaryResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
}

// /users/{id}/profile - the user's id and "about me" type information.
type ProfileResponse struct {
	ID int64 `json:"id"`
	// Name       string `json:"name"`
	// PictureURL string `json:"picture_url"`
	Age int    `json:"age"`
	Bio string `json:"bio"`
	// MaxRadius       float64                `json:"max_radius"`
	// InteractionMode domain.InteractionMode `json:"interaction_mode"`
	// Activities      []ActivityResponse     `json:"activities"`
	// Lat      float64 `json:"lat"`
	// Lon      float64 `json:"lon"`
	// IsOnline bool `json:"is_online"`
}

// Activily responce DTO
type UserActivityResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`          // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
	Experience    int    `json:"experience"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int    `json:"interest_level"` // 1-5 levels: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

// /users/{id}/bio
type UserBioResponse struct {
	ID              int64                  `json:"id"`
	MaxRadius       float64                `json:"max_radius"`
	InteractionMode int                    `json:"interaction_mode"`
	Activities      []UserActivityResponse `json:"activities"`
	Lat             float64                `json:"lat"` // Latitude from the browser API
	Lon             float64                `json:"lon"` // Longitude from the browser API
}

type ActivityResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// /connections
type ConnectionIDResponse struct {
	ID int64 `json:"id"`
}

type ConnectionResponse struct {
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
}

type ConnectionRequest struct {
	ToUserID int64 `json:"to_user_id" validate:"required,number"`
}

type UpdateConnectionRequest struct {
	Status string `json:"status" validate:"required,oneof=accepted declined"`
}

// /chats
type DirectChatRequest struct {
	TargetUserID int64 `json:"target_user_id" validate:"required,number,gt=0"`
}

type ChatResponse struct {
	ID        int64               `json:"id"`
	UserOneID int64               `json:"user_one_id"`
	UserTwoID int64               `json:"user_two_id"`
	CreatedAt time.Time           `json:"created_at"`
	UserOne   UserSummaryResponse `json:"user_one"`
	UserTwo   UserSummaryResponse `json:"user_two"`
}

type MessageResponse struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsViewed  bool      `json:"is_viewed"`
}

