package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/room"
)

type RoomCreateRequest struct {
	// The name of the room to create. This does not need to be globally unique.
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
//	@Security	ApiKey
//	@Success	200 {object} RoomResponse
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/rooms [post]
func (s *Server) handleRoomCreate() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "RoomCreate"))

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

		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			logger.Error("failed to lookup apikey in request context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		opts := room.CreateRoomOpts{Name: request.Name, UserId: userId}
		room, err := s.RoomService.Create(r.Context(), opts)
		if err != nil {
			logger.Error("failed to create room", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := RoomResponse{Id: room.Id, Name: room.Name}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleRoomList lists rooms
//
//	@Summary	Get all rooms
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Security	ApiKey
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

		response := make([]RoomResponse, 0)
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
//	@Security	ApiKey
//	@Success	200
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/rooms/{id} [get]
func (s *Server) handleRoomGet() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "RoomGet"))

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
			logger.Error("failed to encode json response", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleRoomDelete deletes a room
//
//	@Summary	Deletes a room
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"id to delete"
//	@Security	ApiKey
//	@Success	200
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/rooms/{id} [delete]
func (s *Server) handleRoomDelete() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "RoomDelete"))

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			logger.Error("failed to lookup apikey in request context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err = s.RoomService.Delete(r.Context(), room.DeleteRoomOpts{Id: int64(id), UserId: userId}); err != nil {
			switch {
			case errors.Is(err, room.ErrUnauthorized):
				w.WriteHeader(http.StatusUnauthorized)
			case errors.Is(err, room.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
			default:
				logger.Error("failed to delete room", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		return
	}
}

func (c RoomCreateRequest) Validate() bool {
	return c.Name != ""
}
