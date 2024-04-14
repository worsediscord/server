package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/worsediscord/server/services/auth"
	"github.com/worsediscord/server/services/user"
)

var alphaNumericRegex *regexp.Regexp

// swagger:model CreateUserRequest
type CreateUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// swagger:response UserResponse
type UserResponse struct {
	// The globally unique username of the user
	Username string `json:"username"`

	// The nickname of the user
	Nickname string `json:"nickname"`
}

// swagger:response UsersResponse
type UsersResponse []UserResponse

// swagger:model UserLoginRequest
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// swagger:response UserLoginResponse
type UserLoginResponse struct {
	Token string `json:"token"`
}

func init() {
	alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_.]*$")
}

// swagger:route POST /users users createUser
// # Create a user
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//	Responses:
//	  200:
//	  400: Error
//	  409: Error
//	  500: Error
func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: request.Username}); err == nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		opts := user.CreateUserOpts{Username: request.Username, Password: request.Password}
		if err := s.UserService.Create(r.Context(), opts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// swagger:route GET /users users listUsers
// # Gets all users
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//	Responses:
//	  200: UsersResponse
//	  401: Error
//	  500: Error
func (s *Server) handleUserList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.UserService.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var response UsersResponse
		for i := range users {
			response = append(response, UserResponse{Username: users[i].Username, Nickname: users[i].Nickname})
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// swagger:route GET /users/{id} users getUserById
// # Gets a user
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
//	  200: UserResponse
//	  401: Error
func (s *Server) handleUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		u, err := s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: id})
		switch {
		case errors.Is(err, user.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := UserResponse{
			Username: u.Username,
			Nickname: u.Nickname,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// swagger:route GET /users/login users loginUser
// # Logs in a user
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//	Parameters:
//	+ name: credentials
//	  in: body
//	  description: username and password to authenticate with
//	  required: true
//	  type: UserLoginRequest
//	Responses:
//	  200: UserLoginResponse
//	  401: Error
func (s *Server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UserLoginRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		storedUser, err := s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: request.Username})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if request.Password != storedUser.Password {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key := auth.NewApiKey(24, time.Hour*1, storedUser)
		if err = s.AuthService.RegisterKey(key.Token(), key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := UserLoginResponse{Token: key.Token()}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

func (c CreateUserRequest) Validate() bool {
	if c.Username == "" {
		return false
	}

	if !alphaNumericRegex.MatchString(c.Username) {
		return false
	}

	if c.Password == "" {
		return false
	}

	return true
}
