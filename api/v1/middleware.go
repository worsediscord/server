package v2

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/eolso/chat/memcache"
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
func BasicAuthMiddleware(realm string, userDoc *memcache.Document) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, realm)
				return
			}

			//req := NewManagerRequest(MethodList, "", nil)
			//v, err, _ := InterfaceError(req.SendReceive(r.Context(), userReqChan))
			//if err != nil {
			//	basicAuthFailed(w, realm)
			//	return
			//}

			//vUsers, ok := v.([]interface{})
			//if !ok {
			//	basicAuthFailed(w, realm)
			//	return
			//}
			userItems := userDoc.GetAll()
			userMap := make(map[string][]byte)
			for _, item := range userItems {
				//u, ok := vUser.(User)
				var u User
				err := item.Decode(&u)
				if err != nil {
					basicAuthFailed(w, realm)
					return
				}
				password, _ := base64.StdEncoding.DecodeString(u.Password)
				userMap[u.ID] = password
			}

			if subtle.ConstantTimeCompare([]byte(pass), userMap[user]) != 1 {
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
