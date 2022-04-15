package v2

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

type CreateRoomBody struct {
	Name string `json:"name"`
}

type SendMessageBody struct {
	Message string `json:"message"`
}

// CreateRoomHandler creates a room.
func CreateRoomHandler(um *UserManager, rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	// TODO do some permission checking against an API key
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var crb CreateRoomBody
		err = json.Unmarshal(b, &crb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		v := r.Context().Value("userID")
		userID, ok := v.(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		room := NewRoom(crb.Name).WithOwner(userID)
		err = rm.AddRoom(room)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO I don't think we _need_ UserManager here
		err = room.AddUser(um.ActiveUsers[userID])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err = json.Marshal(room)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

// GetRoomHandler gets a room if it is available to the user.
func GetRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO do some permission checking against an API key
		room, err := rm.GetRoom(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		b, err := json.Marshal(room)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

// ListRoomHandler lists all rooms available to the authenticated user via the api key.
func ListRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO do some permission checking against an API key

		response := struct {
			Rooms []*Room `json:"rooms"`
		}{}

		response.Rooms = rm.ListRooms()

		b, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

// DeleteRoomHandler deletes a room.
func DeleteRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO do some permission checking against an API key
		err := rm.DeleteRoom(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ListMessagesHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO do some permission checking against an API key
		room, err := rm.GetRoom(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		response := struct {
			Messages []*Message `json:"messages"`
		}{}

		response.Messages = room.ListMessages()

		b, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func SendMessageHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var crb SendMessageBody
		err = json.Unmarshal(b, &crb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO do some permission checking against an API key
		room, err := rm.GetRoom(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		v := r.Context().Value("userID")
		userID, ok := v.(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, ok := room.GetUser(userID)
		if !ok {
			// TODO this case might be one where we should kick off a cold storage lookup
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		message := NewMessage(string(b), user)
		err = room.SendMessage(message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO really really should have better error reporting/logging
			return
		}

		b, err = json.Marshal(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
