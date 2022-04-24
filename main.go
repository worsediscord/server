package main

import (
	v1 "github.com/eolso/chat/api/v1"
	"github.com/eolso/chat/memcache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const realm = "worsediscord"

func main() {
	var cleaner Cleaner
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		cleaner.Execute()

		os.Exit(0)
	}()

	datastore, err := memcache.Open("state")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create datastore")
	}
	cleaner.RegisterFunc(datastore.Close, "state")

	roomsDoc := datastore.Document("rooms")
	usersDoc := datastore.Document("users")

	roomUsersDoc := datastore.Document("roomUsers")
	userRoomsDoc := datastore.Document("userRooms")

	akm := v1.NewApiKeyManager()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/user", func(r chi.Router) {
			r.Post("/", v1.CreateUserHandler(usersDoc)) // POST /api/v1/user
			r.Get("/", v1.ListUserHandler(usersDoc))    // GET /api/v1/user
			r.Route("/login", func(r chi.Router) {
				r.Use(v1.BasicAuthMiddleware(realm, usersDoc))
				r.Post("/", v1.LoginUserHandler(akm)) // POST /api/v1/user/login
			})
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(v1.ApiAuthMiddleware(akm))
				r.Get("/", v1.GetUserHandler(usersDoc))                        // GET /api/v1/user/{userID}
				r.Delete("/", v1.DeleteUserHandler(usersDoc))                  // DELETE /api/v1/user/{userID}
				r.Get("/room", v1.UserListRoomHandler(userRoomsDoc, roomsDoc)) // GET /api/v1/user/{userID}/room
			})
		})
		// Room routes
		r.Route("/room", func(r chi.Router) {
			r.Use(v1.ApiAuthMiddleware(akm))
			r.Get("/", v1.ListRoomHandler(roomsDoc))                  // GET /api/v1/room
			r.Post("/", v1.CreateRoomHandler(roomsDoc, roomUsersDoc)) // POST /api/v1/room TODO also fix this cuz probs shite

			r.Route("/{roomID}", func(r chi.Router) {
				r.Get("/", v1.GetRoomHandler(roomsDoc))                    // GET /api/v1/room/{roomID}
				r.Delete("/", v1.DeleteRoomHandler(roomsDoc))              // DELETE /api/v1/room/{roomID}
				r.Get("/join", v1.JoinRoomHandler(roomsDoc, roomUsersDoc)) // GET /api/v1/room/{roomID}/join
				r.Post("/invite", v1.InviteRoomHandler(roomsDoc))          // POST /api/v1/room/{roomID}/invite

				r.Route("/message", func(r chi.Router) {
					r.Get("/", v1.ListMessagesHandler(roomsDoc, roomUsersDoc)) // GET /api/v1/room/{roomID}/message
					r.Post("/", v1.SendMessageHandler(roomsDoc, roomUsersDoc)) // POST /api/v1/room/{roomID}/message
				})

				//r.Route("/member", func(r chi.Router) {
				//	r.Patch("/{userID}", v1.PatchRoomUserHandler(rm))
				//})
			})
		})

	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		//cancel()
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
