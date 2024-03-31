package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

type Room struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Messages []Message `json:"-"`
}

func CreateRoomHandler(store storage.Writer[string, Room]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var room Room

		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !room.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		room.Id = base64.StdEncoding.EncodeToString([]byte(room.Name))

		err := store.Write(room.Id, room)
		switch {
		case errors.Is(err, storage.ErrConflict):
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListRoomHandler(store storage.Reader[string, Room]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := store.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(rooms); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetRoomHandler(store storage.Reader[string, Room]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		room, err := store.Read(id)
		switch {
		case errors.Is(err, storage.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(room); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func (r Room) Validate() bool {
	if r.Name == "" {
		return false
	}

	return true
}
