package main

import (
	"context"
	v0 "github.com/eolso/chat/api/v0"
	v2 "github.com/eolso/chat/api/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"path/filepath"
	"time"
)

var ApikeyHeader = "apikey"

var keys = []string{"garysux"}

var creds = map[string]string{
	"glarity": "isbadatrocketleague",
	"eric":    "beesarecute",
}

const realm = "worsediscord"

func main() {
	var s v0.State

	ctx, cancel := context.WithCancel(context.Background())

	roomFlusher, err := v2.NewFileFlusher(filepath.Join("state", "room"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create room flusher")
	}

	userFlusher, err := v2.NewFileFlusher(filepath.Join("state", "user"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create room flusher")
	}

	rm := v2.NewRoomManager().WithFlusher(ctx, roomFlusher)
	um := v2.NewUserManager().WithFlusher(ctx, userFlusher)
	akm := v2.NewApiKeyManager()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v0/messages", func(r chi.Router) {
		r.Use(middleware.BasicAuth("chat", creds))
		r.Get("/", v0.GetMessageHandler(&s))
		r.Post("/", v0.SendMessageHandler(&s))
	})

	// v2 api routes
	r.Route("/api/v2", func(r chi.Router) {
		// user routes
		r.Route("/user", func(r chi.Router) {
			r.Post("/", v2.CreateUserHandler(um)) // POST /api/v2/user
			r.Get("/", v2.ListUserHandler(um))    // GET /api/v2/user TODO this should probably be under an admin route
			r.Route("/login", func(r chi.Router) {
				r.Use(v2.BasicAuthMiddleware(realm, um))
				r.Post("/", v2.LoginUserHandler(um, akm)) // TODO reminder this is shite
			})
		})

		// room routes
		r.Route("/room", func(r chi.Router) {
			r.Use(v2.ApiAuthMiddleware(akm))
			r.Get("/", v2.ListRoomHandler(rm))        // GET /api/v2/room
			r.Post("/", v2.CreateRoomHandler(um, rm)) // POST /api/v2/room TODO also fix this cuz probs shite

			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", v2.GetRoomHandler(rm))       // GET /api/v2/room/{ID}
				r.Delete("/", v2.DeleteRoomHandler(rm)) // DELETE /api/v2/room/{ID}
				r.Get("/message", v2.ListMessagesHandler(rm))
				r.Post("/message", v2.SendMessageHandler(rm))
			})
		})
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
