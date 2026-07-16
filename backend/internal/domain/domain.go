package domain

import (
	"context"
	"time"
)

// Domain:
// 	- Recommendation board
//	- Business language:
//		- find something (find a friend, find a ride, find a group of people)
//		- change profile (update profile), add image, update image
//		- hide profile, add friend, delete friend, ban user
//	- Bounded Context:
//		- Recommendation service:
//			- Profile: publish, hide, update (update preferences)
// 			- Users: seaker, recommendation
// 		- Chat service:
// 			- Messages: send, get, archive, update?, delete?
// 			- Users: sender, receiver
// 		- Account service:
// 			- Account: create (register), login, update, delete
//		- Search...
//		- ...

// Domain:
// - Recommendation Service <- userRepo interface <- UserRepoDB struct

// TODO: update repo methods
type UserRepository interface {
	ProfileByID(ctx context.Context, userID int64) (*Profile, error)
	UpdateProfile(ctx context.Context, profile *Profile) error // maybe use int64 instead of the pointer on an userProfile?
	Recommendations(ctx context.Context, userID int64, limit int) ([]*Profile, error)
}

type ConnectionStatus string

const (
	StatusPending  ConnectionStatus = "pending"
	StatusAccepted ConnectionStatus = "accepted"
	StatusRejected ConnectionStatus = "rejected"
)

type User struct {
	ID           int64  // Postgres SERIAL/BIGSERIAL
	Email        string // unique, private
	PasswordHash string // bcrypt + salt
	CreatedAt    time.Time
}

type Profile struct {
	UserID     int64      `json:"id"`
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	PictureURL string     `json:"picture_url"` // ТЗ: "placeholder image should be shown if no picture"
	Bio        string     `json:"bio"`         // Bio
	Activities []Activity `json:"activities"`
	MaxRadius  float64    `json:"max_radius"`
	Lat        float64    `json:"lat"` // Latitude from the browser API
	Lon        float64    `json:"lon"` // Longitude from the browser API
	IsOnline   bool       `json:"is_online"`
}

type Activity struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`          // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
	Experience    int    `json:"experience"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int    `json:"interest_level"` // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
}

type Connection struct {
	FromUserID int64            `json:"from_user_id"`
	ToUserID   int64            `json:"to_user_id"`
	Status     ConnectionStatus `json:"status"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type Chat struct {
	ID        int64     `json:"id"`
	UserOneID int64     `json:"user_one_id"`
	UserTwoID int64     `json:"user_two_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsViewed  bool      `json:"is_viewed"`           // optional, extra
	IsTyping  bool      `json:"is_typing,omitempty"` // optional, extra
}
