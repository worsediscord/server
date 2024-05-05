package room

import "context"

type Service interface {
	Create(context.Context, CreateRoomOpts) error
	GetRoomById(context.Context, GetRoomByIdOpts) (*Room, error)
	List(context.Context) ([]*Room, error)
	Delete(context.Context, DeleteRoomOpts) error

	Join(context.Context, JoinRoomOpts) error
}
