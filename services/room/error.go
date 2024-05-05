package room

import "errors"

var (
	ErrNotFound     = errors.New("no room found")
	ErrUnauthorized = errors.New("operation is not authorized")
)
