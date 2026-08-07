package domain

// Profile Update structs
type ProfileUpdateParams struct {
	Name            *string
	Age             *int64
	PictureURL      *string
	Bio             *string
	MaxRadius       *float64
	InteractionMode *InteractionMode
	Activities      *[]ActivityInput
	Lat             *float64
	Lon             *float64
}

type ActivityInput struct {
	ActivityID    int64
	Experience    ExperienceLevel
	InterestLevel InterestLevel
}
