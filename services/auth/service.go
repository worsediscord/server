package auth

type Service interface {
	RegisterKey(string, ApiKey) error
	RetrieveKey(string) (ApiKey, error)
	RevokeKey(string) error
}
