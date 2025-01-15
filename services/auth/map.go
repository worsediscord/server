package auth

import (
	"time"

	"github.com/eolso/threadsafe"
)

type Map struct {
	data *threadsafe.Map[string, ApiKey]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, ApiKey](),
	}
}

func (m *Map) RegisterKey(s string, key ApiKey) error {
	m.data.Set(s, key)

	go func() {
		time.Sleep(time.Until(key.ExpiresAt()))
		m.data.Delete(s)
	}()

	return nil
}

func (m *Map) RetrieveKey(s string) (ApiKey, error) {
	key, ok := m.data.Get(s)
	if !ok {
		return ApiKey{}, ErrNotFound
	}

	return key, nil
}

func (m *Map) RevokeKey(s string) error {
	m.data.Delete(s)
	return nil
}
