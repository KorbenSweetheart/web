package storage

import (
	"errors"
)

var (
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidPassFormat     = errors.New("invalid password format")
	ErrEmailIsTaken          = errors.New("email is taken")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCreds          = errors.New("invalid credentials")
	ErrInvalidUsernameFormat = errors.New("invalid username format")
	// ErrUsernameIsTaken       = errors.New("username is taken")
	// ErrSessionNotFound       = errors.New("session not found")
	// ErrInvalidSession        = errors.New("invalid or expired session")
)
