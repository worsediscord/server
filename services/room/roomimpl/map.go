package roomimpl

import (
	"context"

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

	m.data.Set(id, &room.Room{Name: opts.Name, Id: id})

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

func (m *Map) Delete() {
	//TODO implement me
	panic("implement me")
}
