package user

import "errors"

var (
	ErrNotFound        = errors.New("no user found")
	ErrConflict        = errors.New("user already exists")
	ErrInvalidUsername = errors.New("username is invalid")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
)
