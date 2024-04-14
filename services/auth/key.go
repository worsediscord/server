package auth

import (
	"crypto/rand"
	"time"
)

type ApiKey struct {
	payload   any
	token     string
	expiresAt time.Time
}

func NewApiKey(len int, d time.Duration, v any) ApiKey {
	return ApiKey{
		payload:   v,
		token:     string(randBytes(len)),
		expiresAt: time.Now().Add(d),
	}
}

func (a ApiKey) Payload() any {
	return a.payload
}

func (a ApiKey) Token() string {
	return a.token
}

func (a ApiKey) ExpiresAt() time.Time {
	return a.expiresAt
}

func randBytes(length int) []byte {
	const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)

	for index, rbyte := range bytes {
		bytes[index] = validChars[rbyte%byte(len(validChars))]
	}

	return bytes
}
