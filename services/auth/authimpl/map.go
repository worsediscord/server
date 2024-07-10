package authimpl

import (
	"time"

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

	go func() {
		time.Sleep(time.Until(key.ExpiresAt()))
		m.data.Delete(s)
	}()

	return nil
}

func (m *Map) RetrieveKey(s string) (auth.ApiKey, error) {
	key, ok := m.data.Get(s)
	if !ok {
		return auth.ApiKey{}, auth.ErrNotFound
	}

	return key, nil
}

func (m *Map) RevokeKey(s string) error {
	m.data.Delete(s)
	return nil
}
