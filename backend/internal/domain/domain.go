package domain

// Domain:
// 	- Recommendation board
//	- Business language:
//		- find something (find a friend, find a ride, find a group of people)
//		- change profile (update profile), add image, update image
//		- hide profile, add friend, delete friend, ban user
//	- Bounded Context:
//		- Recommendation board:
//			- Profile: publish, update, hide, delete, update preferences
// 			- Users: seaker, recommendation
// 		- Chat:
// 			- Messages: send, get, update?, delete?, archive
// 			- Users: sender, receiver
// 		- Profile...
//		- Search...

// Domain:
// - Recommendation Service <- userRepo interface <- UserRepoDB struct

type Account struct {
	ID    string // maybe use UUID?
	Email string
	User  UserProfile
}

type UserProfile struct {
	ID         string // maybe use UUID?
	FirstName  string
	SecondName string
	Location   map[string]int64 // latitude and longitude
	LookingFor map[string]string
	CanOffer   map[string]string // Maybe just go with 1 entity "LookingFor" and match based on this.
}
