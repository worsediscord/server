package fake

import (
	"github.com/worsediscord/server/services/auth"
)

type AuthService struct {
	ExpectedRegisterKeyError error

	ExpectedRetrieveKeyApiKey auth.ApiKey
	ExpectedRetrieveKeyError  error

	ExpectedRevokeKeyError error
}

func (f *AuthService) RegisterKey(_ string, _ auth.ApiKey) error {
	return f.ExpectedRegisterKeyError
}

func (f *AuthService) RetrieveKey(_ string) (auth.ApiKey, error) {
	return f.ExpectedRetrieveKeyApiKey, f.ExpectedRetrieveKeyError
}

func (f *AuthService) RevokeKey(_ string) error {
	return f.ExpectedRevokeKeyError
}
