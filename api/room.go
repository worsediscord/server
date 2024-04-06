package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

type Room struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Messages []Message `json:"-"`
}

var roomCounter int64

func CreateRoomHandler(store storage.ReadWriter[int64, Room]) http.HandlerFunc {
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

		room.Id = 1000000000000 + roomCounter
		roomCounter += 1

		if _, err := store.Read(room.Id); err == nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := store.Write(room.Id, room); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListRoomHandler(store storage.Reader[int64, Room]) http.HandlerFunc {
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

func GetRoomHandler(store storage.Reader[int64, Room]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		room, err := store.Read(int64(id))
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
