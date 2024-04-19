// Package api worsediscord server API.
//
// HTTP API for interacting with a worsediscord server.
//
//	Schemes: https
//	BasePath: /api
//	Version: 0.0.1
//	Host: test.beesarecute.com
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- api_key:
//	- basic_auth:
//
//	SecurityDefinitions:
//	api_key:
//	  type: apiKey
//	  name: Authorization
//	  in: header
//	basic_auth:
//	  type: basic
//
// swagger:meta
package api

import (
	"log/slog"
	"net/http"

	"github.com/worsediscord/server/services/auth"
	"github.com/worsediscord/server/services/message"
	"github.com/worsediscord/server/services/room"
	"github.com/worsediscord/server/services/user"
)

type Server struct {
	UserService    user.Service
	RoomService    room.Service
	MessageService message.Service
	AuthService    auth.Service

	mux        *http.ServeMux
	logHandler slog.Handler
	middleware []Middleware
}

func NewServer(
	userService user.Service,
	roomService room.Service,
	messageService message.Service,
	authService auth.Service,
	logHandler slog.Handler,
	middleware ...Middleware,
) *Server {
	s := Server{
		UserService:    userService,
		RoomService:    roomService,
		MessageService: messageService,
		AuthService:    authService,
		logHandler:     logHandler,
		mux:            http.NewServeMux(),
		middleware:     middleware,
	}

	authHandler := SessionAuthMiddleware(authService)

	s.mux.Handle("GET /api/users", authHandler(s.handleUserList()))
	s.mux.Handle("POST /api/users", s.handleUserCreate())
	s.mux.Handle("POST /api/users/login", s.handleUserLogin())
	s.mux.Handle("GET /api/users/{id}", authHandler(s.handleUserGet()))

	s.mux.Handle("GET /api/rooms", authHandler(s.handleRoomList()))
	s.mux.Handle("POST /api/rooms", authHandler(s.handleRoomCreate()))
	s.mux.Handle("GET /api/rooms/{id}", authHandler(s.handleRoomGet()))

	s.mux.Handle("GET /api/rooms/{id}/messages", authHandler(s.handleMessageList()))
	s.mux.Handle("POST /api/rooms/{id}/messages", authHandler(s.handleMessageCreate()))

	return &s
}

func (s *Server) AddMiddleware(middleware ...Middleware) {
	s.middleware = append(s.middleware, middleware...)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var h http.Handler

	h = s.mux

	for i := len(s.middleware) - 1; i >= 0; i-- {
		h = s.middleware[i](h)
	}

	h.ServeHTTP(writer, request)
}
