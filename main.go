package main

import (
	"crypto/subtle"
	"fmt"
	v1 "github.com/eolso/chat/api/v1"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var ApikeyHeader = "apikey"

var keys = []string{"garysux"}

var creds = map[string]string{
	"glarity": "isbadatrocketleague",
	"eric":    "beesarecute",
}

func main() {
	r := buildRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func buildRouter() *chi.Mux {
	rm := v1.NewRoomManager("state")
	um := v1.NewUserManager("state")
	urChan, _ := um.Serve()

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

		os.Exit(0)
	}()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v1/rooms", func(r chi.Router) {
		r.Use(middleware.BasicAuth("chat", creds))
		r.Get("/", v1.ListRoomsHandler(rm))         // GET /api/v1/rooms
		r.Post("/", v1.CreateRoomHandler(rm))       // POST /api/v1/rooms
		r.Get("/{ID}", v1.GetRoomHandler(rm))       // GET /api/v1/rooms/{ID}
		r.Delete("/{ID}", v1.RemoveRoomHandler(rm)) // DELETE /api/v1/rooms/{ID}
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/create", v1.CreateUserHandler(urChan)) // POST /api/v1/users/create
	})

	return r
}

func Apikey(keys []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("apikey")
			if key == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var found bool
			for _, k := range keys {
				if subtle.ConstantTimeCompare([]byte(k), []byte(key)) == 1 {
					found = true
				}
			}

			if !found {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
