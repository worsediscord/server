package main

import (
	"github.com/eolso/chat/api/v1"
	"github.com/eolso/memcache"

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

	roomCollection := datastore.Collection("rooms")
	userCollection := datastore.Collection("users")

	//roomsDoc := datastore.Document("rooms")
	//usersDoc := datastore.Document("users")

	//roomUsersDoc := datastore.Document("roomUsers")
	//userRoomsDoc := datastore.Document("userRooms")

	akm := v1.NewApiKeyManager()
	//adminCreds := make(map[string]string)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(nil))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/user", func(r chi.Router) {
			r.Post("/", v1.CreateUserHandler(userCollection)) // POST /api/v1/user
			r.Get("/", v1.ListUserHandler(userCollection))    // GET /api/v1/user
			r.Route("/login", func(r chi.Router) {
				r.Use(v1.BasicAuthMiddleware(realm, userCollection))
				r.Post("/", v1.LoginUserHandler(akm)) // POST /api/v1/user/login
			})
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(v1.ApiAuthMiddleware(akm))
				r.Get("/", v1.GetUserHandler(userCollection)) // GET /api/v1/user/{userID}
				//r.Delete("/", v1.DeleteUserHandler(userCollection))                    // DELETE /api/v1/user/{userID}
				//r.Get("/room", v1.UserListRoomHandler(userCollection, roomCollection)) // GET /api/v1/user/{userID}/room
			})
		})

		// Room routes
		r.Route("/room", func(r chi.Router) {
			r.Use(v1.ApiAuthMiddleware(akm))
			r.Post("/", v1.CreateRoomHandler(roomCollection)) // POST /api/v1/room

			r.Route("/{roomID}", func(r chi.Router) {
				r.Get("/", v1.GetRoomHandler(roomCollection))       // GET /api/v1/room/{roomID}
				r.Delete("/", v1.DeleteRoomHandler(roomCollection)) // DELETE /api/v1/room/{roomID}
				r.Get("/join", v1.JoinRoomHandler(roomCollection))  // GET /api/v1/room/{roomID}/join
				//r.Post("/invite", v1.InviteRoomHandler(roomsDoc))   // POST /api/v1/room/{roomID}/invite

				r.Route("/message", func(r chi.Router) {
					r.Get("/", v1.ListMessagesHandler(roomCollection)) // GET /api/v1/room/{roomID}/message
					r.Post("/", v1.SendMessageHandler(roomCollection)) // POST /api/v1/room/{roomID}/message
				})

				//r.Route("/member", func(r chi.Router) {
				//	r.Patch("/{userID}", v1.PatchRoomUserHandler(rm))
				//})
			})
		})

		// Admin routes
		//r.Route("/admin", func(r chi.Router) {
		//	r.Use(middleware.BasicAuth("admin", adminCreds))
		//	r.Get("/room", v1.AdminListRoomHandler(roomCollection)) // GET /api/v1/room
		//})

	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		//cancel()
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
