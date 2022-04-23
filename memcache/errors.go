package memcache

import (
	"errors"
	"reflect"
)

var ErrKeyNotFound = errors.New("memcache: key not found")
var ErrEmptyItem = errors.New("memcache: Decode(empty item")
var ErrInvalidPath = errors.New("memcache: path must be a directory")

type ErrInvalidDecode struct {
	Type reflect.Type
}

func (e ErrInvalidDecode) Error() string {
	if e.Type == nil {
		return "memcache: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "memcache: Decode(non-pointer " + e.Type.String() + ")"
	}

	return "memcache: Decode(nil " + e.Type.String() + ")"
}
