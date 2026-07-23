package domain

import (
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

type ConnectionStatus string

const (
	StatusPending  ConnectionStatus = "pending"
	StatusAccepted ConnectionStatus = "accepted"
	StatusRejected ConnectionStatus = "rejected"
)

type Account struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"` // Postgres SERIAL/BIGSERIAL
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`  // unique, private
	PasswordHash string    `gorm:"not null" json:"-"`                  // bcrypt + salt
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Profile      Profile   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"profile,omitzero"`
}

type Profile struct {
	UserID            int64             `gorm:"primaryKey" json:"id"`
	Name              string            `gorm:"type:varchar(255);not null" json:"name"`
	Age               int64             `gorm:"column:age" json:"age"`
	PictureURL        string            `gorm:"type:text;default:https://placehold.net/avatar.svg" json:"picture_url"` // Note: "placeholder image should be shown if no picture"
	Bio               string            `gorm:"type:text" json:"bio"`                                                  // Bio
	InteractionModeID int64             `gorm:"default:4" json:"interaction_mode_id"`
	InteractionMode   InteractionMode   `gorm:"foreignKey:InteractionModeID;references:ID;constraint:OnDelete:SET NULL" json:"interaction_mode,omitzero"` // To Preload Activity id and name. "Silent", "Social", "Someone Special / Dating" "Open to anything / Don't care" Mode
	Activities        []ProfileActivity `gorm:"foreignKey:ProfileUserID;references:UserID;constraint:OnDelete:CASCADE" json:"activities"`
	MaxRadius         float64           `gorm:"default:10" json:"max_radius"`
	Lat               float64           `gorm:"-" json:"lat"` // Latitude from the browser API
	Lon               float64           `gorm:"-" json:"lon"` // Longitude from the browser API
	IsOnline          bool              `gorm:"-" json:"is_online"`
}

// InteractionMode directory
type InteractionMode struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"uniqueIndex;not null" json:"title"` // e.g.: "Silent", "Social", "Someone Special / Dating" "Open to anything / Don't care" Mode...
}

// Activily directory
type Activity struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"uniqueIndex;not null" json:"title"` // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
}

// Particular activity in connection to user
type ProfileActivity struct {
	ProfileUserID int64    `gorm:"primaryKey" json:"profile_user_id"`
	ActivityID    int64    `gorm:"primaryKey" json:"activity_id"`
	Experience    int      `gorm:"default:1" json:"experience"`     // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel int      `gorm:"default:3" json:"interest_level"` // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
	Profile       Profile  `gorm:"foreignKey:ProfileUserID;references:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Activity      Activity `gorm:"foreignKey:ActivityID;references:ID;constraint:OnDelete:CASCADE" json:"activity,omitzero"` // To Preload Activity id and name
}

// Provides GORM exact table name
// https://gorm.io/docs/conventions.html#TableName
func (ProfileActivity) TableName() string {
	return "profile_activities"
}

type Connection struct {
	FromUserID int64            `gorm:"primaryKey" json:"from_user_id"`
	ToUserID   int64            `gorm:"primaryKey" json:"to_user_id"`
	Status     ConnectionStatus `gorm:"type:varchar(50);not null" json:"status"`
	UpdatedAt  time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	FromUser   Profile          `gorm:"foreignKey:FromUserID;references:UserID;constraint:OnDelete:CASCADE" json:"-"`
	ToUser     Profile          `gorm:"foreignKey:ToUserID;references:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type Chat struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserOneID int64     `gorm:"not null;index" json:"user_one_id"`
	UserTwoID int64     `gorm:"not null;index" json:"user_two_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UserOne   Profile   `gorm:"foreignKey:UserOneID;references:UserID;constraint:OnDelete:CASCADE" json:"user_one,omitzero"`
	UserTwo   Profile   `gorm:"foreignKey:UserTwoID;references:UserID;constraint:OnDelete:CASCADE" json:"user_two,omitzero"`
	Messages  []Message `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE" json:"messages,omitzero"`
}

type Message struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    int64     `gorm:"not null;index" json:"chat_id"`
	SenderID  int64     `gorm:"not null" json:"sender_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	IsViewed  bool      `gorm:"default:false" json:"is_viewed"` // optional, extra
	IsTyping  bool      `gorm:"-" json:"is_typing,omitempty"`   // optional, extra, e.g., {"type": "typing", "chat_id": 12, "user_id": 42}.
	Chat      Chat      `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Sender    Profile   `gorm:"foreignKey:SenderID;references:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
