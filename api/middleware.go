package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/worsediscord/server/services/auth"
)

type Middleware func(http.Handler) http.Handler

// writeWrapper implements http.ResponseWriter and records a few extra data points.
type writeWrapper struct {
	statusCode   int
	bytesWritten int
	w            http.ResponseWriter
	wroteHeader  bool
}

// RequestLoggerMiddleware logs incoming requests and the status of the response.
func RequestLoggerMiddleware(logHandler slog.Handler, level slog.Level) func(next http.Handler) http.Handler {
	logger := slog.New(logHandler).With(slog.String("middleware", "RequestLoggerMiddleware"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := &writeWrapper{w: w}

			remoteAddr := r.RemoteAddr
			if v := r.Header.Get("CF-Connecting-IP"); v != "" {
				remoteAddr = v
			}

			startTime := time.Now()
			defer func() {
				logger.Log(context.Background(), level,
					fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.Proto),
					slog.String("remote_address", remoteAddr),
					slog.Int("status_code", ww.Status()),
					slog.Int("bytes_written", ww.bytesWritten),
					slog.String("duration", time.Since(startTime).Round(time.Nanosecond).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func SessionAuthMiddleware(authService auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token := r.Header.Get("x-api-key")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			key, err := authService.RetrieveKey(token)
			if err != nil || time.Now().After(key.ExpiresAt()) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "userID", key.Payload())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (w *writeWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *writeWrapper) Write(bytes []byte) (int, error) {
	var err error

	if !w.wroteHeader {
		w.statusCode = http.StatusOK
	}

	w.bytesWritten, err = w.w.Write(bytes)

	return w.bytesWritten, err
}

func (w *writeWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
	w.wroteHeader = true
}

func (w *writeWrapper) Status() int {
	if !w.wroteHeader {
		return http.StatusOK
	}

	return w.statusCode
}
