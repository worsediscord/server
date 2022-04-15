package v2

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
)

func ApiAuthMiddleware(akm *ApiKeyManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key := r.Header.Get("Authorization")
			if key == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID := akm.LookupUser(key)
			if userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "userID", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// BasicAuth implements a simple middleware handler for adding basic http auth to a route.
func BasicAuthMiddleware(realm string, um *UserManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, realm)
				return
			}

			credPass, credUserOk := um.CredMap()[user]
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				basicAuthFailed(w, realm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
