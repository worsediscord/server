package user

import "errors"

var (
	ErrNotFound = errors.New("no user found")
	ErrConflict = errors.New("user already exists")
)
