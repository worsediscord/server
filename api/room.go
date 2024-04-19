package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/room"
)

// swagger:model RoomCreateRequest
type RoomCreateRequest struct {
	Name string `json:"name"`
}

// swagger:response RoomResponse
type RoomResponse struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// swagger:response RoomsResponse
type RoomsResponse []RoomResponse

// swagger:route POST /rooms rooms createRoom
// # Create a room
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//	Parameters:
//	+ name: room
//	  in: body
//	  description: room data
//	  required: true
//	  type: RoomCreateRequest
//	Responses:
//	  200:
//	  400: Error
//	  401: Error
//	  500: Error
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

// swagger:route GET /rooms rooms listRooms
// # Gets all rooms
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//	Responses:
//	  200: RoomsResponse
//	  401: Error
//	  500: Error
func (s *Server) handleRoomList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := s.RoomService.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response RoomsResponse
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

// swagger:route GET /rooms/{id} rooms getRoomById
// # Gets a room
//
//	Produces:
//	- application/json
//	Parameters:
//	+ name: id
//	  in: path
//	  description: id to fetch
//	  required: true
//	  type: string
//	Responses:
//	  200: RoomResponse
//	  401: Error
//	  404: Error
//	  500: Error
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

func (c RoomCreateRequest) Validate() bool {
	return c.Name != ""
}
