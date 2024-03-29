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

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", api.ListUserHandler(userStore))
			r.Post("/", api.CreateUserHandler(userStore))

			r.Get("/{id}", api.GetUserHandler(userStore))
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", api.ListRoomHandler(roomStore))
			r.Post("/", api.CreateRoomHandler(roomStore))

			r.Get("/{id}", api.GetRoomHandler(roomStore))
		})
	})

	panic(http.ListenAndServe(":8069", r))
}
