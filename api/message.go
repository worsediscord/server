package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/worsediscord/server/storage"
)

type Message struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	RoomId    string    `json:"roomId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func CreateMessageHandler(
	messageStore storage.Writer[string, Message],
	roomStore storage.Reader[string, Room],
	userStore storage.Reader[string, User],
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "id")

		// Verify the room exists
		// TODO make a middleware do this
		if _, err := roomStore.Read(roomId); err != nil && errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var message Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message.RoomId = roomId

		// Verify the user exists
		if _, err := userStore.Read(message.UserId); err != nil && errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message.Id = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", roomId, time.Now().UnixNano())))
		message.Timestamp = time.Now()

		if err := messageStore.Write(message.Id, message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

func ListMessageHandler(store storage.Reader[string, Message]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "id")

		allMessages, err := store.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var messages []Message
		for _, message := range allMessages {
			if message.RoomId == roomId {
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
