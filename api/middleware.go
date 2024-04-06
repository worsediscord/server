package api

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"github.com/worsediscord/server/storage"
)

type ApiKeyProperties struct {
	UID       string
	ExpiresAt time.Time
	cancel    context.CancelFunc
}

func NewApiKey(len int, uid string, d time.Duration) (string, ApiKeyProperties) {
	key := string(randBytes(len))

	properties := ApiKeyProperties{
		UID:       uid,
		ExpiresAt: time.Now().Add(d),
	}

	return key, properties
}

func SessionAuthMiddleware(keyStore storage.Reader[string, ApiKeyProperties]) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key := r.Header.Get("Authorization")
			if key == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			keyProperties, err := keyStore.Read(key)
			if err != nil && errors.Is(err, storage.ErrNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "userID", keyProperties.UID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
