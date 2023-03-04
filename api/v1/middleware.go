package v1

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/eolso/memcache"
	"github.com/go-chi/chi"
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

			keyProperties, ok := akm.RetrieveKey(key)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "userID", keyProperties.UID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// BasicAuth implements a simple middleware handler for adding basic http auth to a route.
func BasicAuthMiddleware(realm string, usersCollection *memcache.Collection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				fmt.Println("1")
				basicAuthFailed(w, realm)
				return
			}

			thisUserCollection, ok := usersCollection.GetCollection(user)
			if !ok {
				fmt.Println("2")
				basicAuthFailed(w, realm)
				return
			}

			userDataIface, ok := thisUserCollection.Document(userMetadataKey).Get("_")
			if !ok {
				fmt.Println("3")
				basicAuthFailed(w, realm)
				return
			}
			userData, ok := userDataIface.(*User)
			if !ok {
				fmt.Println("4")
				basicAuthFailed(w, realm)
				return
			}

			userPass, _ := base64.StdEncoding.DecodeString(userData.Metadata.Password)
			if subtle.ConstantTimeCompare([]byte(pass), userPass) != 1 {
				fmt.Println("5")
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

func DocumentMiddleware(collection *memcache.Collection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id := chi.URLParam(r, "ID")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			doc := collection.Document("id")

			switch r.Method {
			case http.MethodGet:
				ctx = context.WithValue(ctx, "document", doc)
			case http.MethodPost:
			case http.MethodPatch:
			case http.MethodDelete:

			}

			next.ServeHTTP(w, r)
		})
	}
}
