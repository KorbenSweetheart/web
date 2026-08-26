package domain

import "errors"

var (
	ErrInvalidEmailFormat      = errors.New("invalid email format")
	ErrInvalidPassFormat       = errors.New("invalid password format")
	ErrEmailIsTaken            = errors.New("email is taken")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidCreds            = errors.New("invalid credentials")
	ErrInvalidUsernameFormat   = errors.New("invalid username format")
	ErrInvalidOrExpiredToken   = errors.New("invalid or expired token")
	ErrIncompleteProfile       = errors.New("incomplete profile")
	ErrConnectionNotFound      = errors.New("connection not found")
	ErrConnectionAlreadyExists = errors.New("connection already exists")
	ErrConnectionNotAllowed    = errors.New("connection not allowed")
	ErrChatNotFound            = errors.New("chat not found")
	ErrNotChatParticipant      = errors.New("user is not a participant in this chat")
	ErrUsersNotConnected       = errors.New("users do not have an accepted connection")
	ErrNoPermissionViewProfile = errors.New("no permission to view profile")
	// ErrUsernameIsTaken          = errors.New("username is taken")
	// ErrInvalidSession           = errors.New("invalid or expired session")
)
