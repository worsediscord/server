package user

import (
	"context"

	"github.com/eolso/threadsafe"
)

type Map struct {
	data *threadsafe.Map[string, *User]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, *User](),
	}
}

func (m *Map) Create(_ context.Context, opts CreateUserOpts) error {
	if _, ok := m.data.Get(opts.Username); ok {
		return ErrConflict
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	u := User{
		Username: opts.Username,
		Nickname: opts.Username,
		Password: opts.Password,
	}

	m.data.Set(opts.Username, &u)

	return nil
}

func (m *Map) GetUserById(_ context.Context, opts GetUserByIdOpts) (*User, error) {
	u, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, ErrNotFound
	}

	return u, nil
}

func (m *Map) List(_ context.Context) ([]*User, error) {
	return m.data.Values(), nil
}

func (m *Map) Delete(_ context.Context, opts DeleteUserOpts) error {
	m.data.Delete(opts.Id)
	return nil
}
