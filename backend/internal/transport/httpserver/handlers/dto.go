package handlers

import "time"

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
	ID                   int64          `json:"id" validate:"required,id"`
	Name                 string         `json:"name"`
	Age                  int64          `json:"age,omitzero"`
	PictureURL           string         `json:"picture_url"`
	Bio                  string         `json:"bio,omitzero"`
	MaxRadius            float64        `json:"max_radius"`
	InteractionModeID    int64          `json:"interaction_mode_id,omitzero"`
	InteractionModeTitle string         `json:"interaction_mode_title,omitzero"`
	Activities           []ActivityResp `json:"activities,omitzero"`
	Lat                  float64        `json:"lat,omitzero"` // Latitude from the browser API
	Lon                  float64        `json:"lon,omitzero"` // Longitude from the browser API
	IsOnline             bool           `json:"is_online,omitzero"`
}

// Activily responce DTO
type ActivityResp struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`          // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
	Experience    int    `json:"experience"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int    `json:"interest_level"` // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

// /users/{id}/bio
type UserBioResponse struct {
	ID        int64    `json:"id" validate:"required,id"`
	Interests []string `json:"interests"`
}

type UpdateProfileRequest struct {
	Name       string  `json:"name" validate:"required"`
	Age        int64   `json:"age" validate:"required,gte=5"`
	PictureURL string  `json:"picture_url" validate:"required,url"`
	Bio        string  `json:"bio" validate:"required"`
	MaxRadius  float64 `json:"max_radius" validate:"required,gt=0"`
}
