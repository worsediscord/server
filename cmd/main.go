package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/cors"
	"github.com/worsediscord/server/api"
	"github.com/worsediscord/server/services/auth/authimpl"
	"github.com/worsediscord/server/services/message/messageimpl"
	"github.com/worsediscord/server/services/room/roomimpl"
	"github.com/worsediscord/server/services/user/userimpl"
)

func main() {
	userService := userimpl.NewMap()
	roomService := roomimpl.NewMap()
	messageService := messageimpl.NewMap()
	authService := authimpl.NewMap()

	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-api-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	s := api.NewServer(userService, roomService, messageService, authService, h, api.RequestLoggerMiddleware(h), corsHandler)

	cmd := NewRootCmd(NewStartCmd(s))

	if err := cmd.Parse(nil); err != nil {
		panic(err)
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
