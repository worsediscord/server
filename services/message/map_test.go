package message

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewMap(t *testing.T) {
	if NewMap() == nil {
		t.Fatal("constructor returned nil")
	}
}

func TestMap_Create(t *testing.T) {
	m := NewMap()

	tests := map[string]struct {
		opts        CreateMessageOpts
		expectedErr error
	}{
		"valid": {
			opts: CreateMessageOpts{
				UserId:  "spiderman",
				RoomId:  100000000000,
				Content: "pizza time",
			},
			expectedErr: nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := m.Create(nil, input.opts)
			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}
		})
	}
}

func TestMap_GetMessageById(t *testing.T) {
	m := NewMap()

	msg, err := m.Create(nil, CreateMessageOpts{UserId: "spiderman", RoomId: 100000000000, Content: "pizza time"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts            GetMessageByIdOpts
		expectedMessage *Message
		expectedErr     error
	}{
		"valid": {
			opts:            GetMessageByIdOpts{Id: msg.Id},
			expectedMessage: msg,
			expectedErr:     nil,
		},
		"not found": {
			opts:            GetMessageByIdOpts{Id: ""},
			expectedMessage: nil,
			expectedErr:     ErrNotFound,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := m.GetMessageById(nil, input.opts)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(u, input.expectedMessage) {
				t.Fatalf("got message %#v, expected %#v", u, input.expectedMessage)
			}
		})
	}
}

func TestMap_List(t *testing.T) {
	nonEmptyMap := NewMap()
	emptyMap := NewMap()

	msg, err := nonEmptyMap.Create(nil, CreateMessageOpts{UserId: "spiderman", RoomId: 100000000000, Content: "pizza time"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m                *Map
		expectedMessages []*Message
		expectedErr      error
	}{
		"non-empty": {
			m:                nonEmptyMap,
			expectedMessages: []*Message{msg},
			expectedErr:      nil,
		},
		"empty": {
			m:                emptyMap,
			expectedMessages: []*Message{},
			expectedErr:      nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			messages, err := input.m.List(nil, ListMessageOpts{})

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(messages, input.expectedMessages) {
				t.Fatalf("got messages %#v, expected %#v", messages, input.expectedMessages)
			}
		})
	}

}
