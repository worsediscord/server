package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

var validCharacters *regexp.Regexp

type User struct {
	Username string `json:"username"`
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

func init() {
	validCharacters = regexp.MustCompile("^[a-zA-Z0-9_.]*$")
}

func CreateUserHandler(store storage.ReadWriter[string, User]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !user.Validate() {
			fmt.Println("ahh")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.Nick = user.Username

		if _, err := store.Read(user.Username); err == nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := store.Write(user.Username, user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListUserHandler(store storage.Reader[string, User]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(users); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetUserHandler(store storage.Reader[string, User]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "id")

		user, err := store.Read(username)
		switch {
		case errors.Is(err, storage.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func LoginUserHandler(userStore storage.Reader[string, User], keyStore storage.Writer[string, ApiKeyProperties]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !user.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		storedUser, err := userStore.Read(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if user.Password != storedUser.Password {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key, props := NewApiKey(24, storedUser.Username, time.Hour*1)
		if err = keyStore.Write(key, props); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type response struct {
			Token string `json:"token"`
		}

		resp := response{Token: key}
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func (u User) Validate() bool {
	if u.Username == "" {
		return false
	}

	if !validCharacters.MatchString(u.Username) {
		return false
	}

	if u.Password == "" {
		return false
	}

	return true
}
