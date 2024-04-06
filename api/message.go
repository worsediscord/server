package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

type Message struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	RoomId    int64  `json:"room_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

func CreateMessageHandler(
	messageStore storage.Writer[string, Message],
	roomStore storage.Reader[int64, Room],
	userStore storage.Reader[string, User],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIdStr := chi.URLParam(r, "id")
		roomId, err := strconv.Atoi(roomIdStr)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify the room exists
		// TODO make a middleware do this
		if _, err := roomStore.Read(int64(roomId)); err != nil && errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var message Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message.RoomId = int64(roomId)

		// Verify the user exists
		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := userStore.Read(userId); err != nil && errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message.UserId = userId
		message.Id = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", roomId, time.Now().UnixNano())))
		message.Timestamp = time.Now().UnixMilli()

		if err := messageStore.Write(message.Id, message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListMessageHandler(store storage.Reader[string, Message]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIdStr := chi.URLParam(r, "id")
		roomId, err := strconv.Atoi(roomIdStr)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		allMessages, err := store.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var messages []Message
		for _, message := range allMessages {
			if message.RoomId == int64(roomId) {
				messages = append(messages, message)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(messages); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		return
	}
}
