package messageimpl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/worsediscord/server/services/message"
)

func TestNewMap(t *testing.T) {
	if NewMap() == nil {
		t.Fatal("constructor returned nil")
	}
}

func TestMap_Create(t *testing.T) {
	m := NewMap()

	tests := map[string]struct {
		opts        message.CreateMessageOpts
		expectedErr error
	}{
		"valid": {
			opts: message.CreateMessageOpts{
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

	msg, err := m.Create(nil, message.CreateMessageOpts{UserId: "spiderman", RoomId: 100000000000, Content: "pizza time"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts            message.GetMessageByIdOpts
		expectedMessage *message.Message
		expectedErr     error
	}{
		"valid": {
			opts:            message.GetMessageByIdOpts{Id: msg.Id},
			expectedMessage: msg,
			expectedErr:     nil,
		},
		"not found": {
			opts:            message.GetMessageByIdOpts{Id: ""},
			expectedMessage: nil,
			expectedErr:     message.ErrNotFound,
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

	msg, err := nonEmptyMap.Create(nil, message.CreateMessageOpts{UserId: "spiderman", RoomId: 100000000000, Content: "pizza time"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m                *Map
		expectedMessages []*message.Message
		expectedErr      error
	}{
		"non-empty": {
			m:                nonEmptyMap,
			expectedMessages: []*message.Message{msg},
			expectedErr:      nil,
		},
		"empty": {
			m:                emptyMap,
			expectedMessages: []*message.Message{},
			expectedErr:      nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			messages, err := input.m.List(nil, message.ListMessageOpts{})

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(messages, input.expectedMessages) {
				t.Fatalf("got messages %#v, expected %#v", messages, input.expectedMessages)
			}
		})
	}

}
