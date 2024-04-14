package user

import "context"

type Service interface {
	Create(context.Context, CreateUserOpts) error
	GetUserById(context.Context, GetUserByIdOpts) (*User, error)
	List(context.Context) ([]*User, error)
	Delete()
}
