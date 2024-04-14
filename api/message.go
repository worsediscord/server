package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/message"
	"github.com/worsediscord/server/services/room"
	"github.com/worsediscord/server/services/user"
)

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type ListMessageResponse []message.Message

func (s *Server) handleMessageCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify the room exists
		if _, err = s.RoomService.GetRoomById(r.Context(), room.GetRoomByIdOpts{Id: int64(roomId)}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var request CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify the user exists
		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: userId}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		opts := message.CreateMessageOpts{
			UserId:  userId,
			RoomId:  int64(roomId),
			Content: request.Content,
		}

		if err = s.MessageService.Create(r.Context(), opts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

func (s *Server) handleMessageList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		allMessages, err := s.MessageService.List(r.Context(), message.ListMessageOpts{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response ListMessageResponse
		for _, msg := range allMessages {
			if msg.RoomId == int64(roomId) {
				response = append(response, *msg)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}
