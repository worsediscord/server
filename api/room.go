package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/room"
)

func (s *Server) handleRoomCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request room.CreateRoomOpts

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.RoomService.Create(r.Context(), request); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO return id probably

		return
	}
}

func (s *Server) handleRoomList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := s.RoomService.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(rooms); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

func (s *Server) handleRoomGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response, err := s.RoomService.GetRoomById(r.Context(), room.GetRoomByIdOpts{Id: int64(id)})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}
