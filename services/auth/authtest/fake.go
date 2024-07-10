package authtest

import "github.com/worsediscord/server/services/auth"

type FakeAuthService struct {
	ExpectedRegisterKeyError error

	ExpectedRetrieveKeyApiKey auth.ApiKey
	ExpectedRetrieveKeyError  error

	ExpectedRevokeKeyError error
}

func (f *FakeAuthService) RegisterKey(_ string, _ auth.ApiKey) error {
	return f.ExpectedRegisterKeyError
}

func (f *FakeAuthService) RetrieveKey(_ string) (auth.ApiKey, error) {
	return f.ExpectedRetrieveKeyApiKey, f.ExpectedRetrieveKeyError
}

func (f *FakeAuthService) RevokeKey(_ string) error {
	return f.ExpectedRevokeKeyError
}
