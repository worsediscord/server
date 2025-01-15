package auth

import (
	"reflect"
	"testing"
	"time"
)

func TestMap_RegisterKey(t *testing.T) {
	m := NewMap()

	if err := m.RegisterKey("key", NewApiKey(1, time.Second, nil)); err != nil {
		t.Fatal(err)
	}
}

func TestMap_RetrieveKey(t *testing.T) {
	m := NewMap()

	tests := map[string]struct {
		key      string
		length   int
		duration time.Duration
		payload  any
	}{
		"valid": {
			key:      "key",
			length:   8,
			duration: time.Minute,
			payload:  "hello",
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if err := m.RegisterKey(input.key, NewApiKey(input.length, input.duration, input.payload)); err != nil {
				t.Fatal(err)
			}

			key, err := m.RetrieveKey(input.key)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(key.Payload(), input.payload) {
				t.Fatalf("got %v, expected %v", key.Payload(), input.payload)
			}
		})
	}
}

func TestMap_RevokeKey(t *testing.T) {
	m := NewMap()

	// Currently this always returns nil, but might be worth having a key not found error at some point
	if err := m.RevokeKey("key"); err != nil {
		t.Fatal(err)
	}
}
