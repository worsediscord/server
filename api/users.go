package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUserHandler(store storage.Writer[string, User]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := store.Write(user.Name, user)
		switch {
		case errors.Is(err, storage.ErrConflict):
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListUserHandler(store storage.Reader[string, User]) func(w http.ResponseWriter, r *http.Request) {
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

func GetUserHandler(store storage.Reader[string, User]) func(w http.ResponseWriter, r *http.Request) {
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
