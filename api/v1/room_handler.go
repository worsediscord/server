package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolso/chat/memcache"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

type CreateRoomBody struct {
	Name string `json:"name"`
}

type InviteUserBody struct {
	ID string `json:"id"`
}

type SendMessageBody struct {
	Message string `json:"message"`
}

type PatchRoomUserBody struct {
	DisplayName string `json:"display_name"`
}

// CreateRoomHandler creates a room.
func CreateRoomHandler(roomDoc memcache.DocumentWriter, roomUserMapDoc memcache.DocumentWriter) func(w http.ResponseWriter, r *http.Request) {
	// TODO do some permission checking against an API key
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var crb CreateRoomBody
		if err = json.Unmarshal(b, &crb); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		room := NewRoom(crb.Name).WithOwner(userID)
		if err = roomDoc.Set(room.ID, room); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Update the room -> user map with the new owner
		var roomUsers []Identity
		roomUsers = append(roomUsers, Identity{ID: userID, Name: userID})
		if err = roomUserMapDoc.Set(room.ID, roomUsers); err != nil {
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
func GetRoomHandler(roomDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO do some permission checking against an API key
		var room Room
		err := roomDoc.Get(roomID).Decode(&room)
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
func ListRoomHandler(roomDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Rooms []interface{} `json:"rooms"`
		}{}

		roomItems := roomDoc.GetAll()
		for _, item := range roomItems {
			var room Room
			if err := item.Decode(&room); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.Rooms = append(response.Rooms, room)
		}

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
func DeleteRoomHandler(roomDoc memcache.DocumentWriter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		roomDoc.Delete(roomID)

		w.WriteHeader(http.StatusOK)
	}
}

func JoinRoomHandler(roomDoc memcache.DocumentReadWriter, roomUserMap memcache.DocumentReadWriter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Fetch the authed userID from the api key context
		userID, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Verify that the room exists, and fetch data needed for later
		var room Room
		if err := roomDoc.Get(roomID).Decode(&room); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify that the user doesn't already exist in the room
		var roomUserList []Identity
		if err := roomUserMap.Get(roomID).Decode(&roomUserList); err != nil {
			// This implies that the room doesn't exist and shouldn't really ever occur
			if !errors.Is(err, memcache.ErrEmptyItem) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		for _, roomUser := range roomUserList {
			if userID == roomUser.ID {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		if !room.IsInvited(userID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		roomUserList = append(roomUserList, Identity{ID: userID, Name: userID})

		if err := roomUserMap.Set(roomID, roomUserList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := room.SendSystemMessage(fmt.Sprintf("%s has joined the room. Say hi!", userID)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := roomDoc.Set(roomID, room); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func InviteRoomHandler(roomDoc memcache.DocumentReadWriter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var iub InviteUserBody
		if err = json.Unmarshal(b, &iub); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify that the room exists, and fetch data needed for later
		var room Room
		if err = roomDoc.Get(roomID).Decode(&room); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		room.InviteUser(iub.ID)
		if err = roomDoc.Set(roomID, room); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ListMessagesHandler(roomDoc memcache.DocumentReader, roomUserDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var room Room
		err := roomDoc.Get(roomID).Decode(&room)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		userID, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check that the user exists in the room
		var roomUsers []Identity
		if err = roomUserDoc.Get(roomID).Decode(&roomUsers); err != nil {
			// This implies that the room user mapping doesn't exist but the room does. This shouldn't really ever occur.
			// TODO possibly try and recover from this state.
			if !errors.Is(err, memcache.ErrEmptyItem) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		var found bool
		for _, ru := range roomUsers {
			if userID == ru.ID {
				found = true
			}
		}
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response := struct {
			Messages []Message `json:"messages"`
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

func SendMessageHandler(roomDoc memcache.DocumentReadWriter, roomUserDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var smb SendMessageBody
		if err = json.Unmarshal(b, &smb); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify that the room exists, and fetch data needed for later
		var room Room
		if err = roomDoc.Get(roomID).Decode(&room); err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
			return
		}

		userID, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check that the user exists in the room
		var roomUsers []Identity
		if err = roomUserDoc.Get(roomID).Decode(&roomUsers); err != nil {
			// This implies that the room user mapping doesn't exist but the room does. This shouldn't really ever occur.
			// TODO possibly try and recover from this state.
			if !errors.Is(err, memcache.ErrEmptyItem) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		var roomUser Identity
		var found bool
		for _, ru := range roomUsers {
			if userID == ru.ID {
				roomUser = ru
				found = true
			}
		}

		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		message := NewMessage(smb.Message, roomUser)
		err = room.SendMessage(message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // TODO really really should have better error reporting/logging
			return
		}

		err = roomDoc.Set(roomID, room)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

//
//func PatchRoomUserHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		id := chi.URLParam(r, "ID")
//		if id == "" {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		targetUserID := chi.URLParam(r, "userID")
//		if id == "" {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		b, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		var requestBody PatchRoomUserBody
//		err = json.Unmarshal(b, &requestBody)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		// TODO do some permission checking against an API key
//		room, err := rm.GetRoom(id)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest) // TODO maybe should be more useful here (not found vs unauthorized etc)
//			return
//		}
//
//		v := r.Context().Value("userID")
//		userID, ok := v.(string)
//		if !ok {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		// Shortcut for targeting yourself
//		if targetUserID == "@me" {
//			targetUserID = userID
//		}
//
//		// TODO eventually we should do some ACL checks here. But for now you can only change your own name.
//		if userID != targetUserID {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		err = room.UpdateUserName(targetUserID, requestBody.DisplayName)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest) // TODO really really should have better error reporting/logging
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//	}
//}
