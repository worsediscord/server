package usertest

import (
	"context"

	"github.com/worsediscord/server/services/user"
)

type FakeUserService struct {
	ExpectedUser  *user.User
	ExpectedUsers []*user.User
	ExpectedError error
}

func (f *FakeUserService) Create(_ context.Context, _ user.CreateUserOpts) error {
	return f.ExpectedError
}

func (f *FakeUserService) GetUserById(_ context.Context, _ user.GetUserByIdOpts) (*user.User, error) {
	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) List(_ context.Context) ([]*user.User, error) {
	return f.ExpectedUsers, f.ExpectedError
}

func (f *FakeUserService) Delete(_ context.Context, _ user.DeleteUserOpts) error {
	return f.ExpectedError
}
