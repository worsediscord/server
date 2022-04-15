package v1

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
func LoginUserHandler(akm *ApiKeyManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()

		apikey, properties := NewApiKey(24, user, time.Hour)
		akm.RegisterKey(apikey, properties)

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
