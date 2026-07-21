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
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`       // Postgres SERIAL/BIGSERIAL
	Email        string    `gorm:"uniqueIndex;not null;column:email" json:"email"`     // unique, private
	PasswordHash string    `gorm:"not null;column:password_hash" json:"password_hash"` // bcrypt + salt
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

type Profile struct {
	UserID            int64             `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name              string            `gorm:"type:varchar(255);not null;column:name" json:"name"`
	Age               int               `gorm:"column:age" json:"age"`
	PictureURL        string            `gorm:"type:text;default:https://placehold.net/avatar.svg;column:picture_url" json:"picture_url"` // Note: "placeholder image should be shown if no picture"
	Bio               string            `gorm:"type:text;column:bio" json:"bio"`                                                          // Bio
	InteractionModeID int64             `gorm:"column:interaction_mode_id;default:4" json:"interaction_mode_id"`
	InteractionMode   InteractionMode   `gorm:"foreignKey:InteractionModeID;references:ID" json:"interaction_mode,omitzero"` // To Preload Activity id and name. "Silent", "Social", "Someone Special / Dating" "Open to anything / Don't care" Mode
	Activities        []ProfileActivity `gorm:"foreignKey:ProfileUserID;references:UserID" json:"activities"`
	MaxRadius         float64           `gorm:"default:10;column:max_radius" json:"max_radius"`
	Lat               float64           `gorm:"-" json:"lat"` // Latitude from the browser API
	Lon               float64           `gorm:"-" json:"lon"` // Longitude from the browser API
	IsOnline          bool              `gorm:"-" json:"is_online"`
}

// InteractionMode directory
type InteractionMode struct {
	ID    int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Title string `gorm:"uniqueIndex;not null;column:title" json:"title"` // e.g.: "Silent", "Social", "Someone Special / Dating" "Open to anything / Don't care" Mode...
}

// Activily directory
type Activity struct {
	ID    int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Title string `gorm:"uniqueIndex;not null;column:title" json:"title"` // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
}

// Particular activity in connection to user
type ProfileActivity struct {
	ProfileUserID int64    `gorm:"primaryKey;column:profile_user_id"`
	ActivityID    int64    `gorm:"primaryKey;column:activity_id"`
	Experience    int      `gorm:"column:experience;default:1" json:"experience"`                // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int      `gorm:"column:interest_level;default:3" json:"interest_level"`        // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
	Activity      Activity `gorm:"foreignKey:ActivityID;references:ID" json:"activity,omitzero"` // To Preload Activity id and name
}

// Provides GORM exact table name
// https://gorm.io/docs/conventions.html#TableName
func (ProfileActivity) TableName() string {
	return "profile_activities"
}

type Connection struct {
	FromUserID int64            `gorm:"primaryKey;column:from_user_id" json:"from_user_id"`
	ToUserID   int64            `gorm:"primaryKey;column:to_user_id" json:"to_user_id"`
	Status     ConnectionStatus `gorm:"type:varchar(50);not null;column:status" json:"status"`
	UpdatedAt  time.Time        `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

type Chat struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserOneID int64     `gorm:"not null;index;column:user_one_id" json:"user_one_id"`
	UserTwoID int64     `gorm:"not null;index;column:user_two_id" json:"user_two_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

type Message struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ChatID    int64     `gorm:"not null;index;column:chat_id" json:"chat_id"`
	SenderID  int64     `gorm:"not null;column:sender_id" json:"sender_id"`
	Content   string    `gorm:"type:text;not null;column:content" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	IsViewed  bool      `gorm:"default:false;column:is_viewed" json:"is_viewed"` // optional, extra
	IsTyping  bool      `gorm:"-" json:"is_typing,omitempty"`                    // optional, extra, e.g., {"type": "typing", "chat_id": 12, "user_id": 42}.
}
