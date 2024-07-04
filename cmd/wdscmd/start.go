package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/cors"
	"github.com/worsediscord/server/cmd"
	"github.com/worsediscord/server/services/auth/authimpl"
	"github.com/worsediscord/server/services/message/messageimpl"
	"github.com/worsediscord/server/services/room/roomimpl"
	"github.com/worsediscord/server/services/user/userimpl"
	"github.com/worsediscord/server/util"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/worsediscord/server/api"
)

type StartCmd struct {
	Port string

	LogLevel    string
	LogFormat   string
	LogRequests bool
}

func NewStartCmd() *StartCmd {
	return &StartCmd{
		Port:        "8069",
		LogLevel:    "info",
		LogFormat:   "text",
		LogRequests: false,
	}
}

func (s *StartCmd) Name() string {
	return "start"
}

func (s *StartCmd) Description() string {
	return "Start a worsediscord server"
}

func (s *StartCmd) Parse(args []string) error {
	fs := flag.NewFlagSet(s.Name(), flag.ExitOnError)

	fs.StringVar(&s.Port, "p", s.Port, "TCP Port to listen on.")
	fs.StringVar(&s.Port, "port", s.Port, cmd.LongFlagUsage("p"))

	fs.StringVar(&s.LogLevel, "log-level", s.LogLevel, "log level")
	fs.StringVar(&s.LogFormat, "log-format", s.LogFormat, "log format (text | json | disabled)")
	fs.BoolVar(&s.LogRequests, "log-requests", s.LogRequests, "Enable logging of requests")

	fs.Usage = func() {
		_, _ = fmt.Fprint(flag.CommandLine.Output(), cmd.HelpString(s, fs))
	}

	return fs.Parse(args)
}

func (s *StartCmd) Run() error {
	var logHandler slog.Handler
	var middleware []api.Middleware

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
	middleware = append(middleware, corsHandler)

	switch strings.ToLower(s.LogFormat) {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: util.StringToLogLevel(s.LogLevel)})
	case "disabled":
		logHandler = util.NopLogHandler
	default:
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: util.StringToLogLevel(s.LogLevel)})
	}

	if s.LogRequests {
		middleware = append(middleware, api.RequestLoggerMiddleware(logHandler, slog.LevelDebug))
	}

	server := api.NewServer(userService, roomService, messageService, authService, logHandler, middleware...)

	return http.ListenAndServe(":"+s.Port, server)
}
