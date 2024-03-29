package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/worsediscord/server/api"
	"github.com/worsediscord/server/storage"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userStore := storage.NewMap[string, api.User]()
	roomStore := storage.NewMap[string, api.Room]()
	messageStore := storage.NewMap[string, api.Message]()

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", api.ListUserHandler(userStore))
			r.Post("/", api.CreateUserHandler(userStore))

			r.Get("/{id}", api.GetUserHandler(userStore))
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", api.ListRoomHandler(roomStore))
			r.Post("/", api.CreateRoomHandler(roomStore))

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", api.GetRoomHandler(roomStore))

				r.Get("/messages", api.ListMessageHandler(messageStore))
				r.Post("/messages", api.CreateMessageHandler(messageStore, roomStore, userStore))
			})
		})
	})

	panic(http.ListenAndServe(":8069", r))
}
