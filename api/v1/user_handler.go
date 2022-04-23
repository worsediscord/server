package v2

import (
	"encoding/json"
	"errors"
	"github.com/eolso/chat/memcache"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateUserBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserListRoomResponse struct {
	Name string
	ID   string
}

func CreateUserHandler(userDoc memcache.DocumentWriter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var cub CreateUserBody
		err = json.Unmarshal(b, &cub)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := NewUser(cub.Name).WithPassword(cub.Password)
		err = userDoc.Set(cub.Name, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err = json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func GetUserHandler(userDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var u User
		err := userDoc.Get(userID).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func ListUserHandler(userDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO do some permission checking against an API key
		response := struct {
			Users []interface{} `json:"users"`
		}{}

		userItems := userDoc.GetAll()
		for _, item := range userItems {
			var user User
			if err := item.Decode(&user); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.Users = append(response.Users, user)
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

func DeleteUserHandler(userDoc memcache.DocumentWriter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userDoc.Delete(userID)

		w.WriteHeader(http.StatusOK)
	}
}

func UserListRoomHandler(userRoomMapDoc memcache.DocumentReadWriter, roomDoc memcache.DocumentReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authedUserID, ok := r.Context().Value("userID").(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if userID == "@me" {
			userID = authedUserID
		}

		var roomIDs []string
		if err := userRoomMapDoc.Get(userID).Decode(&roomIDs); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response []UserListRoomResponse
		var room Room
		var mapUpdated bool
		for i, roomID := range roomIDs {
			if err := roomDoc.Get(roomID).Decode(&room); err != nil {
				// The room no longer exists, delete it.
				if errors.Is(err, memcache.ErrEmptyItem) {
					roomIDs = append(roomIDs[:i], roomIDs[i+1:]...)
					mapUpdated = true
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			response = append(response, UserListRoomResponse{Name: room.Name, ID: room.ID})
		}

		if mapUpdated {
			if err := userRoomMapDoc.Set(userID, roomIDs); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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

// LoginUserHandler uses the basic auth header to authenticate a user then returns an API key to use.
func LoginUserHandler(akm *ApiKeyManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()

		apikey, properties := NewApiKey(24, user, time.Hour)

		err := akm.RegisterKey(apikey, properties)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := struct {
			Apikey string `json:"apikey"`
		}{
			Apikey: apikey,
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
