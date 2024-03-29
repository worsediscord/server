package storage

import "errors"

var (
	ErrNotFound = errors.New("no value found")
	ErrConflict = errors.New("value already exists")
)
