package room

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
		opts         CreateRoomOpts
		expectedRoom *Room
		expectedErr  error
	}{
		"valid": {
			opts:         CreateRoomOpts{Name: "the big apple", UserId: "spidey"},
			expectedRoom: &Room{Id: padding, Name: "the big apple", Users: []string{"spidey"}, Admins: []string{"spidey"}},
			expectedErr:  nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			createdRoom, err := m.Create(nil, input.opts)
			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(createdRoom, input.expectedRoom) {
				t.Fatalf("got %v, expected %v", createdRoom, input.expectedRoom)
			}
		})
	}
}

func TestMap_GetRoomById(t *testing.T) {
	m := NewMap()

	createdRoom, err := m.Create(nil, CreateRoomOpts{Name: "the big apple", UserId: "spiderman"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts         GetRoomByIdOpts
		expectedRoom *Room
		expectedErr  error
	}{
		"valid": {
			opts:         GetRoomByIdOpts{Id: createdRoom.Id},
			expectedRoom: &Room{Id: 100000000000, Name: "the big apple", Users: []string{"spiderman"}, Admins: []string{"spiderman"}},
			expectedErr:  nil,
		},
		"not found": {
			opts:         GetRoomByIdOpts{Id: 1},
			expectedRoom: nil,
			expectedErr:  ErrNotFound,
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

	createdRoom, err := nonEmptyMap.Create(nil, CreateRoomOpts{Name: "the big apple", UserId: "spiderman"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m             *Map
		expectedRooms []*Room
		expectedErr   error
	}{
		"non-empty": {
			m:             nonEmptyMap,
			expectedRooms: []*Room{{Id: createdRoom.Id, Name: "the big apple", Users: []string{"spiderman"}, Admins: []string{"spiderman"}}},
			expectedErr:   nil,
		},
		"empty": {
			m:             emptyMap,
			expectedRooms: []*Room{},
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

	roomToDelete, err := m.Create(nil, CreateRoomOpts{Name: "the big apple", UserId: "spiderman"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	unauthorizedRoom, err := m.Create(nil, CreateRoomOpts{Name: "the big apple (backup)", UserId: "spiderman"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		roomId      int64
		userId      string
		expectedErr error
	}{
		"valid": {
			roomId:      roomToDelete.Id,
			userId:      "spiderman",
			expectedErr: nil,
		},
		"not found": {
			roomId:      1,
			expectedErr: ErrNotFound,
		},
		"unauthorized": {
			roomId:      unauthorizedRoom.Id,
			userId:      "batman",
			expectedErr: ErrUnauthorized,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			err := m.Delete(nil, DeleteRoomOpts{Id: input.roomId, UserId: input.userId})

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}
		})
	}
}

func TestMap_Join(t *testing.T) {
	m := NewMap()

	createdRoom, err := m.Create(nil, CreateRoomOpts{Name: "the big apple", UserId: "spiderman"})
	if err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts        JoinRoomOpts
		expectedErr error
	}{
		"valid": {
			opts:        JoinRoomOpts{Id: createdRoom.Id, UserId: "batman"},
			expectedErr: nil,
		},
		"not found": {
			opts:        JoinRoomOpts{Id: 1, UserId: "batman"},
			expectedErr: ErrNotFound,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			err := m.Join(nil, input.opts)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}
		})
	}
}
