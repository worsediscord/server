package usertest

import (
	"context"

	"github.com/worsediscord/server/services/user"
)

type FakeUserService struct {
	ExpectedCreateError error

	ExpectedGetUserByIdUser  *user.User
	ExpectedGetUserByIdError error

	ExpectedListUsers []*user.User
	ExpectedListError error

	ExpectedDeleteError error
}

func (f *FakeUserService) Create(_ context.Context, _ user.CreateUserOpts) error {
	return f.ExpectedCreateError
}

func (f *FakeUserService) GetUserById(_ context.Context, _ user.GetUserByIdOpts) (*user.User, error) {
	return f.ExpectedGetUserByIdUser, f.ExpectedGetUserByIdError
}

func (f *FakeUserService) List(_ context.Context) ([]*user.User, error) {
	return f.ExpectedListUsers, f.ExpectedListError
}

func (f *FakeUserService) Delete(_ context.Context, _ user.DeleteUserOpts) error {
	return f.ExpectedDeleteError
}
