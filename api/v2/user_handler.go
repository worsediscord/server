package v2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateUserBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUserHandler(um *UserManager) func(w http.ResponseWriter, r *http.Request) {
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

		err = um.AddUser(user)
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

func ListUserHandler(um *UserManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// LoginUserHandler uses the basic auth header to authenticate a user then returns an API key to use.
func LoginUserHandler(um *UserManager, akm *ApiKeyManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()

		apikey := NewApiKey(24, time.Hour)

		// TODO fix all of this shit lol
		uid := um.LookupUser(user)
		if uid == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		akm.RegisterKey(uid, apikey)

		response := struct {
			Apikey string `json:"apikey"`
		}{
			Apikey: akm.RetrieveKey(uid).Key(),
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
