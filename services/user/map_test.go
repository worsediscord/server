package user

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
		opts        CreateUserOpts
		expectedErr error
	}{
		"initial valid": {
			opts:        CreateUserOpts{Username: "spiderman", Password: "uncleben123"},
			expectedErr: nil,
		},
		"duplicate user": {
			opts:        CreateUserOpts{Username: "spiderman", Password: "uncleben123"},
			expectedErr: ErrConflict,
		},
		"invalid user": {
			opts:        CreateUserOpts{Username: "", Password: "uncleben123"},
			expectedErr: ErrInvalidUsername,
		},
		"invalid password": {
			opts:        CreateUserOpts{Username: "spiderman2", Password: "ben"},
			expectedErr: ErrInvalidPassword,
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

func TestMap_GetUserById(t *testing.T) {
	m := NewMap()

	if err := m.Create(nil, CreateUserOpts{Username: "spiderman", Password: "uncleben123"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts         GetUserByIdOpts
		expectedUser *User
		expectedErr  error
	}{
		"valid": {
			opts:         GetUserByIdOpts{Id: "spiderman"},
			expectedUser: &User{Username: "spiderman", Nickname: "spiderman", Password: "uncleben123"},
			expectedErr:  nil,
		},
		"invalid": {
			opts:         GetUserByIdOpts{Id: "antman"},
			expectedUser: nil,
			expectedErr:  ErrNotFound,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := m.GetUserById(nil, input.opts)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(u, input.expectedUser) {
				t.Fatalf("got user %#v, expected %#v", u, input.expectedUser)
			}
		})
	}
}

func TestMap_List(t *testing.T) {
	nonEmptyMap := NewMap()
	emptyMap := NewMap()

	if err := nonEmptyMap.Create(nil, CreateUserOpts{Username: "spiderman", Password: "uncleben123"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m             *Map
		expectedUsers []*User
		expectedErr   error
	}{
		"non-empty": {
			m:             nonEmptyMap,
			expectedUsers: []*User{{Username: "spiderman", Nickname: "spiderman", Password: "uncleben123"}},
			expectedErr:   nil,
		},
		"empty": {
			m:             emptyMap,
			expectedUsers: []*User{},
			expectedErr:   nil,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			users, err := input.m.List(nil)

			if !errors.Is(err, input.expectedErr) {
				t.Fatalf("got error %q, expected %q", err, input.expectedErr)
			}

			if !reflect.DeepEqual(users, input.expectedUsers) {
				t.Fatalf("got users %#v, expected %#v", users, input.expectedUsers)
			}
		})
	}

}

func TestMap_Delete(t *testing.T) {
	m := NewMap()

	if err := m.Delete(nil, DeleteUserOpts{}); err != nil {
		t.Fatalf("got error %q, expected nil", err)
	}
}
