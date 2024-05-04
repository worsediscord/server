package userimpl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/worsediscord/server/services/user"
)

func TestNewMap(t *testing.T) {
	if NewMap() == nil {
		t.Fatal("constructor returned nil")
	}
}

func TestMap_Create(t *testing.T) {
	m := NewMap()

	tests := map[string]struct {
		opts        user.CreateUserOpts
		expectedErr error
	}{
		"initial valid": {
			opts:        user.CreateUserOpts{Username: "spiderman", Password: "uncleben123"},
			expectedErr: nil,
		},
		"duplicate user": {
			opts:        user.CreateUserOpts{Username: "spiderman", Password: "uncleben123"},
			expectedErr: user.ErrConflict,
		},
		"invalid user": {
			opts:        user.CreateUserOpts{Username: "", Password: "uncleben123"},
			expectedErr: user.ErrInvalidUsername,
		},
		"invalid password": {
			opts:        user.CreateUserOpts{Username: "spiderman2", Password: "ben"},
			expectedErr: user.ErrInvalidPassword,
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

	if err := m.Create(nil, user.CreateUserOpts{Username: "spiderman", Password: "uncleben123"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		opts         user.GetUserByIdOpts
		expectedUser *user.User
		expectedErr  error
	}{
		"valid": {
			opts:         user.GetUserByIdOpts{Id: "spiderman"},
			expectedUser: &user.User{Username: "spiderman", Nickname: "spiderman", Password: "uncleben123"},
			expectedErr:  nil,
		},
		"invalid": {
			opts:         user.GetUserByIdOpts{Id: "antman"},
			expectedUser: nil,
			expectedErr:  user.ErrNotFound,
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

	if err := nonEmptyMap.Create(nil, user.CreateUserOpts{Username: "spiderman", Password: "uncleben123"}); err != nil {
		t.Fatalf("failed to prepopulate map: %v", err)
	}

	tests := map[string]struct {
		m             *Map
		expectedUsers []*user.User
		expectedErr   error
	}{
		"non-empty": {
			m:             nonEmptyMap,
			expectedUsers: []*user.User{{Username: "spiderman", Nickname: "spiderman", Password: "uncleben123"}},
			expectedErr:   nil,
		},
		"empty": {
			m:             emptyMap,
			expectedUsers: []*user.User{},
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

	if err := m.Delete(nil, user.DeleteUserOpts{}); err != nil {
		t.Fatalf("got error %q, expected nil", err)
	}
}
