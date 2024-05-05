package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/worsediscord/server/services/auth"
	"github.com/worsediscord/server/services/user"
)

type UserCreateRequest struct {
	// The globally unique username of the user.
	Username string `json:"username,omitempty"`

	// The password to set. Must be at least 8 characters long.
	Password string `json:"password,omitempty"`
}

type UserResponse struct {
	// The globally unique username of the user.
	Username string `json:"username"`

	// The nickname of the user.
	Nickname string `json:"nickname"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

// handleUserCreate creates a user
//
//	@Summary	Create a user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		credentials	body	UserCreateRequest	true	"username and password to create user with"
//	@Success	200
//	@Failure	400
//	@Failure	409
//	@Failure	500
//	@Router		/users [post]
func (s *Server) handleUserCreate() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "UserCreate"))

	return func(w http.ResponseWriter, r *http.Request) {
		var request UserCreateRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Error("failed to decode request", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.Validate() {
			logger.Error("request failed validation")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: request.Username}); err == nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		opts := user.CreateUserOpts{Username: request.Username, Password: request.Password}
		if err := s.UserService.Create(r.Context(), opts); err != nil {
			logger.Error("failed to create user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserList returns a list of users
//
//	@Summary	List users
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	ApiKey
//	@Success	200	{array}	UserResponse
//	@Failure	400
//	@Failure	500
//	@Router		/users [get]
func (s *Server) handleUserList() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "UserList"))

	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.UserService.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var response []UserResponse
		for i := range users {
			response = append(response, UserResponse{Username: users[i].Username, Nickname: users[i].Nickname})
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("failed to encode response", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserGet gets a user
//
//	@Summary	Gets a user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"id to fetch"
//	@Security	ApiKey
//	@Success	200	{object}	UserResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/users/{id} [get]
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

// handleUserLogin logs in a user and returns a token
//
//	@Summary	Logs in a user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	BasicAuth
//	@Success	200	{object}	UserLoginResponse
//	@Failure	401
//	@Failure	500
//	@Router		/users/login [post]
func (s *Server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		storedUser, err := s.UserService.GetUserById(r.Context(), user.GetUserByIdOpts{Id: username})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if password != storedUser.Password {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key := auth.NewApiKey(24, time.Hour*1, storedUser.Username)
		if err = s.AuthService.RegisterKey(key.Token(), key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := UserLoginResponse{Token: key.Token()}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDelete deletes a user
//
//	@Summary	Deletes a user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"id to delete"
//	@Security	ApiKey
//	@Success	200
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/users/{id} [delete]
func (s *Server) handleUserDelete() http.HandlerFunc {
	logger := slog.New(s.logHandler).With(slog.String("handle", "UserDelete"))

	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("userID").(string)
		if !ok {
			logger.Error("failed to lookup apikey in request context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if userId != r.PathValue("id") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := s.UserService.Delete(r.Context(), user.DeleteUserOpts{Id: userId}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO the user's token is pretty much still valid

		return
	}
}

func (c UserCreateRequest) Validate() bool {
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
