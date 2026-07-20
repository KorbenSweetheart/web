package storage

import (
	"errors"
)

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidPassFormat  = errors.New("invalid password format")
	ErrEmailIsTaken       = errors.New("email is taken")
	ErrInvalidCreds       = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	// ErrInvalidUsernameFormat = errors.New("invalid username format")
	// ErrUsernameIsTaken       = errors.New("username is taken")
)
