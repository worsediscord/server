package main

import (
	"context"
	"fmt"
	v0 "github.com/eolso/chat/api/v0"
	v1 "github.com/eolso/chat/api/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const realm = "worsediscord"

func main() {
	var s v0.State

	ctx, cancel := context.WithCancel(context.Background())

	roomFlusher, err := v1.NewFileFlusher(filepath.Join("state", "room"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create room flusher")
	}

	userFlusher, err := v1.NewFileFlusher(filepath.Join("state", "user"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create room flusher")
	}

	rm := v1.NewRoomManager().WithFlusher(ctx, roomFlusher)
	um := v1.NewUserManager().WithFlusher(ctx, userFlusher)
	akm := v1.NewApiKeyManager()

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		err := rm.Flush()
		if err != nil {
			fmt.Println(err)
		}

		err = um.Flush()
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 2) // TODO be better than this

		os.Exit(0)
	}()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v0/messages", func(r chi.Router) {
		var creds = map[string]string{
			"glarity": "isbadatrocketleague",
			"eric":    "beesarecute",
		}
		r.Use(middleware.BasicAuth("chat", creds))
		r.Get("/", v0.GetMessageHandler(&s))
		r.Post("/", v0.SendMessageHandler(&s))
	})

	// v1 api routes
	r.Route("/api/v1", func(r chi.Router) {
		// user routes
		r.Route("/user", func(r chi.Router) {
			r.Post("/", v1.CreateUserHandler(um)) // POST /api/v1/user
			r.Get("/", v1.ListUserHandler(um))    // GET /api/v1/user TODO this should probably be under an admin route
			r.Route("/login", func(r chi.Router) {
				r.Use(v1.BasicAuthMiddleware(realm, um))
				r.Post("/", v1.LoginUserHandler(akm)) // TODO reminder this is shite
			})
		})

		// room routes
		r.Route("/room", func(r chi.Router) {
			r.Use(v1.ApiAuthMiddleware(akm))
			r.Get("/", v1.ListRoomHandler(rm))        // GET /api/v1/room
			r.Post("/", v1.CreateRoomHandler(um, rm)) // POST /api/v1/room TODO also fix this cuz probs shite

			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", v1.GetRoomHandler(um, rm))   // GET /api/v1/room/{ID}
				r.Delete("/", v1.DeleteRoomHandler(rm)) // DELETE /api/v1/room/{ID}
				r.Get("/message", v1.ListMessagesHandler(rm))
				r.Post("/message", v1.SendMessageHandler(rm))
			})
		})
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
