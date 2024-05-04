package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/worsediscord/server/services/auth"
)

type Middleware func(http.Handler) http.Handler

func RequestLoggerMiddleware(h slog.Handler) func(next http.Handler) http.Handler {
	logger := slog.New(h)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteAddr := r.RemoteAddr
			if v := r.Header.Get("CF-Connecting-IP"); v != "" {
				remoteAddr = v
			}

			logger.Info(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
				slog.String("remote_addr", remoteAddr),
			)

			next.ServeHTTP(w, r)
		})
	}
}

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

			ctx = context.WithValue(ctx, "userID", key.Payload())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
