package roomimpl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/worsediscord/server/services/room"
)

func TestNewMap(t *testing.T) {
	if NewMap() == nil {
		t.Fatal("constructor returned nil")
	}
}

func TestMap_Create(t *testing.T) {
	m := NewMap()

	tests := map[string]struct {
		opts        room.CreateRoomOpts
		expectedErr error
	}{
		"valid": {
			opts:        room.CreateRoomOpts{Name: "the big apple"},
			expectedErr: nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if err := m.Create(nil, input.opts); !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}
		})
	}
}

func TestMap_GetRoomById(t *testing.T) {
	m := NewMap()

	if err := m.Create(nil, room.CreateRoomOpts{Name: "the big apple"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts         room.GetRoomByIdOpts
		expectedRoom *room.Room
		expectedErr  error
	}{
		"valid": {
			opts:         room.GetRoomByIdOpts{Id: 100000000000},
			expectedRoom: &room.Room{Id: 100000000000, Name: "the big apple"},
			expectedErr:  nil,
		},
		"not found": {
			opts:         room.GetRoomByIdOpts{Id: 1},
			expectedRoom: nil,
			expectedErr:  room.ErrNotFound,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := m.GetRoomById(nil, input.opts)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(u, input.expectedRoom) {
				t.Fatalf("got room %#v, expected %#v", u, input.expectedRoom)
			}
		})
	}
}

func TestMap_List(t *testing.T) {
	nonEmptyMap := NewMap()
	emptyMap := NewMap()

	if err := nonEmptyMap.Create(nil, room.CreateRoomOpts{Name: "the big apple"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m             *Map
		expectedRooms []*room.Room
		expectedErr   error
	}{
		"non-empty": {
			m:             nonEmptyMap,
			expectedRooms: []*room.Room{{Id: 100000000000, Name: "the big apple"}},
			expectedErr:   nil,
		},
		"empty": {
			m:             emptyMap,
			expectedRooms: []*room.Room{},
			expectedErr:   nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			rooms, err := input.m.List(nil)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(rooms, input.expectedRooms) {
				t.Fatalf("got rooms %#v, expected %#v", rooms, input.expectedRooms)
			}
		})
	}

}

func TestMap_Delete(t *testing.T) {
	m := NewMap()

	if err := m.Delete(nil, room.DeleteRoomOpts{}); err != nil {
		t.Fatalf("got error %q, expected nil", err)
	}
}
