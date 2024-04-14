package userimpl

import (
	"context"

	"github.com/eolso/threadsafe"
	"github.com/worsediscord/server/services/user"
)

type Map struct {
	data *threadsafe.Map[string, *user.User]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, *user.User](),
	}
}

func (m *Map) Create(_ context.Context, opts user.CreateUserOpts) error {
	if _, ok := m.data.Get(opts.Username); ok {
		return user.ErrConflict
	}

	u := user.User{
		Username: opts.Username,
		Nickname: opts.Username,
		Password: opts.Password,
	}

	m.data.Set(opts.Username, &u)

	return nil
}

func (m *Map) GetUserById(_ context.Context, opts user.GetUserByIdOpts) (*user.User, error) {
	u, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (m *Map) List(_ context.Context) ([]*user.User, error) {
	return m.data.Values(), nil
}

func (m *Map) Delete() {
	//TODO implement me
	panic("implement me")
}
