package domain

import (
	"fmt"
	"strings"
	"time"
)

type ConnectionStatus int

const (
	Pending  ConnectionStatus = iota + 1 // 1
	Accepted                             // 2
	Declined                             // 3
)

func (s ConnectionStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Accepted:
		return "accepted"
	case Declined:
		return "declined"
	default:
		return "unknown"
	}
}

func ParseConnectionStatus(s string) (ConnectionStatus, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "pending":
		return Pending, nil
	case "accepted":
		return Accepted, nil
	case "declined":
		return Declined, nil
	default:
		return 0, fmt.Errorf("invalid connection status: %s", s)
	}
}

type InteractionMode int

const (
	OpenToAnything InteractionMode = iota + 1 // 1
	Social                                    // 2
	Silent                                    // 3
	Dating                                    // 4
)

type ExperienceLevel int

const (
	Beginner     ExperienceLevel = iota + 1 // 1
	ActiveNovice                            // 2
	Intermediate                            // 3
	Advanced                                // 4
	Professional                            // 5
)

type InterestLevel int

const (
	NotInterested InterestLevel = iota + 1 // 1
	OpenToIt                               // 2
	Interested                             // 3
	Highly                                 // 4
	ActivelyLook                           // 5
)

type Account struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"` // Postgres SERIAL/BIGSERIAL
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`  // unique, private
	PasswordHash string    `gorm:"not null" json:"-"`                  // bcrypt + salt
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Profile      Profile   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"profile,omitzero"`
}

type RefreshToken struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID int64     `gorm:"index;not null" json:"account_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"` // hashedRefreshToken
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"` // e.g., time.Now().Add(30 * 24 * time.Hour),
	Account   Account   `gorm:"foreignKey:AccountID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type Profile struct {
	UserID     int64  `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`
	PictureURL string `gorm:"type:text;default:https://placehold.net/avatar.svg" json:"picture_url"` // Note: "placeholder image should be shown if no picture"
	Age        int    `gorm:"type:smallint;column:age" json:"age"`
	// Birthday        time.Time         `gorm:"type:date" json:"birthday"`                                             // to get age AGE(birthday) in SQL or time.Since(profile.Birthday) in Go
	Bio             string            `gorm:"type:text" json:"bio"` // Bio
	MaxRadius       float64           `gorm:"default:10" json:"max_radius"`
	InteractionMode InteractionMode   `gorm:"type:smallint;column:interaction_mode;not null;default:1" json:"interaction_mode"` // "Open to anything", "Silent", "Social", "Dating" Mode
	Activities      []ProfileActivity `gorm:"foreignKey:ProfileUserID;references:UserID;constraint:OnDelete:CASCADE" json:"activities"`
	Lat             float64           `gorm:"column:lat" json:"lat"` // Latitude from the browser API
	Lon             float64           `gorm:"column:lon" json:"lon"` // Longitude from the browser API
}

// Activily dictionary
type Activity struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"uniqueIndex;not null" json:"title"` // e.g.: "Running", "Paddle", "Gym", "Cycling", "CrossFit", "Football"...
}

// Particular activity in connection to user
type ProfileActivity struct {
	ProfileUserID int64           `gorm:"primaryKey" json:"profile_user_id"`
	ActivityID    int64           `gorm:"primaryKey" json:"activity_id"`
	Experience    ExperienceLevel `gorm:"type:smallint;column:experience_level;not null;default:1" json:"experience_level"` // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
	InterestLevel InterestLevel   `gorm:"type:smallint;column:interest_level;not null;default:3" json:"interest_level"`     // 1-5 levels: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
	Profile       Profile         `gorm:"foreignKey:ProfileUserID;references:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Activity      Activity        `gorm:"foreignKey:ActivityID;references:ID;constraint:OnDelete:CASCADE" json:"activity,omitzero"` // To Preload Activity id and name
}

// Provides GORM exact table name
// https://gorm.io/docs/conventions.html#TableName
func (ProfileActivity) TableName() string {
	return "profile_activities"
}

// Connections
// When a user sees a recommendation that they find interesting, they can request to connect with them.
// Users must be able to see a list of connection requests, where they can accept or dismiss requests.
// It must be possible to disconnect with a user, if they are no longer interesting.
//
// Profiles are viewable by other users, only if:
// - They are recommended
// - There is an outstanding connection request
// - They are connected
type Connection struct {
	FromUserID int64            `gorm:"primaryKey" json:"from_user_id"`
	ToUserID   int64            `gorm:"primaryKey" json:"to_user_id"`
	Status     ConnectionStatus `gorm:"type:smallint;not null;default:1" json:"status"`
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
