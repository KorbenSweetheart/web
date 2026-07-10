package domain

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

type AcID string // maybe just to test???
type Account struct {
	ID     AcID // maybe use UUID? Implement generation via factory or inside the repo as an interface, just to have a capability to swap it if needed.
	Email  string
	UserID string
}

type UserProfile struct {
	ID          string // maybe use UUID? or maybe skip it?
	UserName    string
	Avatar      string // path to the file
	Bio         string
	Location    Location
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
