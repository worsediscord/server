package fake

import (
	"context"

	"github.com/worsediscord/server/services/user"
)

type UserService struct {
	ExpectedCreateError error

	ExpectedGetUserByIdUser  *user.User
	ExpectedGetUserByIdError error

	ExpectedListUsers []*user.User
	ExpectedListError error

	ExpectedDeleteError error
}

func (f *UserService) Create(_ context.Context, _ user.CreateUserOpts) error {
	return f.ExpectedCreateError
}

func (f *UserService) GetUserById(_ context.Context, _ user.GetUserByIdOpts) (*user.User, error) {
	return f.ExpectedGetUserByIdUser, f.ExpectedGetUserByIdError
}

func (f *UserService) List(_ context.Context) ([]*user.User, error) {
	return f.ExpectedListUsers, f.ExpectedListError
}

func (f *UserService) Delete(_ context.Context, _ user.DeleteUserOpts) error {
	return f.ExpectedDeleteError
}
