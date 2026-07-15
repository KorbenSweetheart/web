package domain

import "context"

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

type UserRepository interface {
	UserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	GetRecommendations(ctx context.Context, userID string, limit int) ([]*User, error)
}

type User struct {
	ID          string // maybe use UUID? or maybe skip it?
	Email       string
	Password    string
	UserName    string
	AvatarURL   string // path to the file
	AboutMe     string
	Location    Location
	Interests   []string // one option of matching
	LookingFor  map[string]string
	CanOffer    map[string]string // Maybe just go with 1 entity "LookingFor" and match based on this.
	FriendsList []string
}

type Location struct {
	Country   string
	City      string
	Latitude  float64
	Longitude float64
}
