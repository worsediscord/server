package storage

import (
	"github.com/eolso/threadsafe"
)

type Map[K comparable, V any] struct {
	tsMap *threadsafe.Map[K, V]
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		tsMap: threadsafe.NewMap[K, V](),
	}
}

func (m Map[K, V]) Read(key K) (V, error) {
	v, ok := m.tsMap.Get(key)
	if !ok {
		return *new(V), ErrNotFound
	}

	return v, nil
}

func (m Map[K, V]) ReadAll() ([]V, error) {
	return m.tsMap.Values(), nil
}

func (m Map[K, V]) Write(key K, v V) error {
	if _, ok := m.tsMap.Get(key); ok {
		return ErrConflict
	}

	m.tsMap.Set(key, v)
	return nil
}
