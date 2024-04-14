package api

import (
	"context"
	"net/http"

	"github.com/worsediscord/server/services/auth"
)

type Middleware func(http.Handler) http.Handler

func SessionAuthMiddleware(authService auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			key, err := authService.RetrieveKey(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "userID", key.Token())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
