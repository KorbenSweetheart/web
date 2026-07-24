package handlers

import "time"

// AuthHandler DTO
type RegisterRequest struct {
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
	Name       string `json:"name" validate:"required,name"`
	PictureURL string `json:"picture_url" validate:"required,picture_url"`
}

// /users/{id}/profile
type ProfileResponse struct {
	UserID     int64   `json:"user_id"`
	Name       string  `json:"name"`
	Age        int64   `json:"age"`
	PictureURL string  `json:"picture_url"`
	Bio        string  `json:"bio"`
	MaxRadius  float64 `json:"max_radius"`
}

// /users/{id}/bio
type UserBioResponse struct {
	UserID    int64    `json:"user_id"`
	Interests []string `json:"interests"`
}

type UpdateProfileRequest struct {
	Name       string  `json:"name" validate:"required"`
	Age        int64   `json:"age" validate:"required,gte=5"`
	PictureURL string  `json:"picture_url" validate:"required,url"`
	Bio        string  `json:"bio" validate:"required"`
	MaxRadius  float64 `json:"max_radius" validate:"required,gt=0"`
}
