package handlers

import (
	"match-me-api/internal/domain"
	"time"
)

// AuthHandler DTO
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserHandler DTO

// /users/{id}
type UserSummaryResponse struct {
	ID         int64  `json:"id" validate:"required,id"`
	Name       string `json:"name" validate:"required,name"`
	PictureURL string `json:"picture_url" validate:"required,picture_url"`
}

// /users/{id}/profile - the user's id and "about me" type information.
type ProfileResponse struct {
	ID              int64                  `json:"id" validate:"required,id"`
	Name            string                 `json:"name"`
	Age             int64                  `json:"age,omitzero"`
	PictureURL      string                 `json:"picture_url"`
	Bio             string                 `json:"bio,omitzero"`
	MaxRadius       float64                `json:"max_radius,omitzero"`
	InteractionMode domain.InteractionMode `json:"interaction_mode,omitzero"`
	Activities      []ActivityResponse     `json:"activities,omitzero"`
	Lat             float64                `json:"lat,omitzero"` // Latitude from the browser API
	Lon             float64                `json:"lon,omitzero"` // Longitude from the browser API
	IsOnline        bool                   `json:"is_online,omitzero"`
}

// Activily responce DTO
type ActivityResponse struct {
	ID            int64                  `json:"id"`
	Title         string                 `json:"title"`          // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
	Experience    domain.ExperienceLevel `json:"experience"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel domain.InterestLevel   `json:"interest_level"` // 1-5 levels: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

// /users/{id}/bio
type UserBioResponse struct {
	ID              int64                  `json:"id" validate:"required,id"`
	MaxRadius       float64                `json:"max_radius"`
	InteractionMode domain.InteractionMode `json:"interaction_mode,omitzero"`
	Activities      []ActivityResponse     `json:"activities,omitzero"`
}

type UpdateProfileRequest struct {
	ID              int64                  `json:"id" validate:"required,id"`
	Name            string                 `json:"name" validate:"min=3"`
	Age             int64                  `json:"age" validate:"gte=5"`
	PictureURL      string                 `json:"picture_url" validate:"url"`
	Bio             string                 `json:"bio"`
	MaxRadius       float64                `json:"max_radius" validate:"gt=0"`
	InteractionMode domain.InteractionMode `json:"interaction_mode,omitzero" validate:"gte=1,lte=4"`
	Activities      []ActivityRequest      `json:"activities,omitzero"`
}

type ActivityRequest struct {
	ID            int64                  `json:"id" validate:"required,id"`
	Title         string                 `json:"title" validate:"min=3"`                // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
	Experience    domain.ExperienceLevel `json:"experience" validate:"gte=1,lte=5"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel domain.InterestLevel   `json:"interest_level" validate:"gte=1,lte=5"` // 1-5 levels: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

type UpdateLocationRequest struct {
	Lat float64 `json:"lat" validate:"required,latitude"`
	Lon float64 `json:"lon" validate:"required,longitude"`
}
