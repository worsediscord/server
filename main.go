package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/worsediscord/server/api"
	"github.com/worsediscord/server/storage"
)

func main() {
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)

	userStore := storage.NewMap[string, api.User]()
	roomStore := storage.NewMap[int64, api.Room]()
	messageStore := storage.NewMap[string, api.Message]()
	keyStore := storage.NewMap[string, api.ApiKeyProperties]()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", api.CreateUserHandler(userStore))
			r.Post("/login", api.LoginUserHandler(userStore, keyStore))

			// Authenticated routes
			r.Route("/", func(r chi.Router) {
				r.Use(api.SessionAuthMiddleware(keyStore))
				r.Get("/", api.ListUserHandler(userStore))
				r.Get("/{id}", api.GetUserHandler(userStore))
			})
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Use(api.SessionAuthMiddleware(keyStore))

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
