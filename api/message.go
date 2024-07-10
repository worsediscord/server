package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/worsediscord/server/services/message"
	"github.com/worsediscord/server/services/room"
	"github.com/worsediscord/server/services/user"
)

type MessageCreateRequest struct {
	// The content of the message.
	Content string `json:"content"`
}

type MessageResponse struct {
	// The unique username of the message author.
	UserId string `json:"user_id,omitempty"`

	// The content of the message.
	Content string `json:"content,omitempty"`

	// Time since epoch in milliseconds.
	Timestamp int64 `json:"timestamp,omitempty"`
}

// handleMessageCreate creates a message
//
//	@Summary	Create a message
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string					true	"room id to create message in"
//	@Param		content	body	MessageCreateRequest	true	"content to create message with"
//	@Security	ApiKey
//	@Success	200
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/rooms/{id}/messages [post]
func (s *Server) handleMessageCreate() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handler", "MessageCreate"))

	return func(w http.ResponseWriter, r *http.Request) {
		var logAttrs []slog.Attr

		roomId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify the room exists
		if _, err = s.RoomService.GetRoomById(r.Context(), room.GetRoomByIdOpts{Id: roomId}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logAttrs = append(logAttrs, slog.Int64("room_id", roomId))

		var request MessageCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify the user exists
		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			logger.LogAttrs(r.Context(), slog.LevelError, "failed to lookup apikey in request context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: userId}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logAttrs = append(logAttrs, slog.String("user_id", userId))

		opts := message.CreateMessageOpts{
			UserId:  userId,
			RoomId:  roomId,
			Content: request.Content,
		}

		if _, err = s.MessageService.Create(r.Context(), opts); err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "failed to create message", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.LogAttrs(r.Context(), slog.LevelInfo, "message created", logAttrs...)

		return
	}
}

// handleMessageList returns a list of messages for a given room
//
//	@Summary	List messages
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"room id to list messages from"
//	@Security	ApiKey
//	@Success	200	{array}	MessageResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/rooms/{id}/messages [get]
func (s *Server) handleMessageList() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handler", "MessageList"))

	return func(w http.ResponseWriter, r *http.Request) {
		var logAttrs []slog.Attr

		roomId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		logAttrs = append(logAttrs, slog.Int64("room_id", roomId))

		// Verify the user exists
		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			logger.LogAttrs(r.Context(), slog.LevelError, "failed to lookup apikey in request context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: userId}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logAttrs = append(logAttrs, slog.String("user_id", userId))

		messages, err := s.MessageService.List(r.Context(), message.ListMessageOpts{RoomId: roomId})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]MessageResponse, 0)
		for _, msg := range messages {
			response = append(response, MessageResponse{
				UserId:    msg.UserId,
				Content:   msg.Content,
				Timestamp: msg.Timestamp,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.LogAttrs(r.Context(), slog.LevelInfo, "messages listed", logAttrs...)

		return
	}
}
