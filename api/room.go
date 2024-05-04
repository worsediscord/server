package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/room"
)

type RoomCreateRequest struct {
	Name string `json:"name"`
}

type RoomResponse struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// handleRoomCreate creates a room
//
//	@Summary	Create a room
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Param		name	body	RoomCreateRequest	true	"room data"
//	@Security	BasicAuth
//	@Success	200
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/rooms [post]
func (s *Server) handleRoomCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RoomCreateRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		opts := room.CreateRoomOpts{Name: request.Name}
		if err := s.RoomService.Create(r.Context(), opts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO return id probably

		return
	}
}

// handleRoomList lists rooms
//
//	@Summary	Get all rooms
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Failure	401
//	@Failure	500
//	@Router		/rooms [get]
func (s *Server) handleRoomList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := s.RoomService.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response []RoomResponse
		for i := range rooms {
			response = append(response, RoomResponse{Id: rooms[i].Id, Name: rooms[i].Name})
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleRoomGet gets a room
//
//	@Summary	Gets a room
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"id to fetch"
//	@Success	200
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/rooms/{id} [get]
func (s *Server) handleRoomGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gotRoom, err := s.RoomService.GetRoomById(r.Context(), room.GetRoomByIdOpts{Id: int64(id)})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(RoomResponse{Id: gotRoom.Id, Name: gotRoom.Name}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

func (s *Server) handleRoomDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err = s.RoomService.Delete(r.Context(), room.DeleteRoomOpts{int64(id)}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		return
	}
}

func (c RoomCreateRequest) Validate() bool {
	return c.Name != ""
}
