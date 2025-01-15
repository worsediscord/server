package room

import (
	"context"
	"slices"

	"github.com/eolso/threadsafe"
)

const padding = 100000000000

type Map struct {
	data        *threadsafe.Map[int64, *Room]
	padding     int64
	roomCounter int64
}

func NewMap() *Map {
	return &Map{
		data:        threadsafe.NewMap[int64, *Room](),
		padding:     padding,
		roomCounter: 0,
	}
}

func (m *Map) Create(_ context.Context, opts CreateRoomOpts) (*Room, error) {
	id := m.padding + m.roomCounter
	r := &Room{Name: opts.Name, Id: id, Users: []string{opts.UserId}, Admins: []string{opts.UserId}}

	m.data.Set(id, r)
	m.roomCounter += 1

	return r, nil
}

func (m *Map) GetRoomById(_ context.Context, opts GetRoomByIdOpts) (*Room, error) {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, ErrNotFound
	}

	return r, nil
}

func (m *Map) List(_ context.Context) ([]*Room, error) {
	return m.data.Values(), nil
}

func (m *Map) Delete(_ context.Context, opts DeleteRoomOpts) error {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return ErrNotFound
	}

	if !opts.Force && !slices.Contains(r.Admins, opts.UserId) {
		return ErrUnauthorized
	}

	m.data.Delete(opts.Id)
	return nil
}

func (m *Map) Join(_ context.Context, opts JoinRoomOpts) error {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return ErrNotFound
	}

	if slices.Contains(r.Users, opts.UserId) {
		return nil
	}

	// This is still pretty vulnerable to race conditions
	r.Users = append(r.Users, opts.UserId)
	m.data.Set(r.Id, r)

	return nil
}
