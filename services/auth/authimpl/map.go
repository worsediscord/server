package authimpl

import (
	"github.com/eolso/threadsafe"
	"github.com/worsediscord/server/services/auth"
)

type Map struct {
	data *threadsafe.Map[string, auth.ApiKey]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, auth.ApiKey](),
	}
}

func (m *Map) RegisterKey(s string, key auth.ApiKey) error {
	m.data.Set(s, key)
	return nil
}

func (m *Map) RetrieveKey(s string) (auth.ApiKey, error) {
	key, ok := m.data.Get(s)
	if !ok {
		return auth.ApiKey{}, auth.ErrNotFound
	}

	return key, nil
}

func (m *Map) RevokeKey(_ string) error {
	//TODO implement me
	panic("implement me")
}
