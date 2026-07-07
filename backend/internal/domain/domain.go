package domain

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
