package roomimpl

import (
	"context"
	"slices"

	"github.com/eolso/threadsafe"
	"github.com/worsediscord/server/services/room"
)

type Map struct {
	data        *threadsafe.Map[int64, *room.Room]
	padding     int64
	roomCounter int64
}

func NewMap() *Map {
	return &Map{
		data:        threadsafe.NewMap[int64, *room.Room](),
		padding:     100000000000,
		roomCounter: 0,
	}
}

func (m *Map) Create(_ context.Context, opts room.CreateRoomOpts) error {
	id := m.padding + m.roomCounter

	m.data.Set(id, &room.Room{Name: opts.Name, Id: id, Users: []string{opts.UserId}, Admins: []string{opts.UserId}})

	m.roomCounter += 1

	return nil
}

func (m *Map) GetRoomById(_ context.Context, opts room.GetRoomByIdOpts) (*room.Room, error) {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, room.ErrNotFound
	}

	return r, nil
}

func (m *Map) List(_ context.Context) ([]*room.Room, error) {
	return m.data.Values(), nil
}

func (m *Map) Delete(_ context.Context, opts room.DeleteRoomOpts) error {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return room.ErrNotFound
	}

	if !opts.Force && !slices.Contains(r.Admins, opts.UserId) {
		return room.ErrUnauthorized
	}

	m.data.Delete(opts.Id)
	return nil
}

func (m *Map) Join(_ context.Context, opts room.JoinRoomOpts) error {
	r, ok := m.data.Get(opts.Id)
	if !ok {
		return room.ErrNotFound
	}

	if slices.Contains(r.Users, opts.UserId) {
		return nil
	}

	// This is still pretty vulnerable to race conditions
	r.Users = append(r.Users, opts.UserId)
	m.data.Set(r.Id, r)

	return nil
}
