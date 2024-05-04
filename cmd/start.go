package main

import (
	"flag"
	"net/http"

	"github.com/worsediscord/server/api"
)

type StartCmd struct {
	server *api.Server
	port   string
}

func NewStartCmd(s *api.Server) *StartCmd {
	return &StartCmd{
		server: s,
		port:   "8069",
	}
}

func (s *StartCmd) Name() string {
	return "start"
}

func (s *StartCmd) Description() string {
	return "Start a worsediscord server"
}

func (s *StartCmd) Parse(args []string) error {
	flags := flag.NewFlagSet(s.Name(), flag.ExitOnError)
	flags.StringVar(&s.port, "p", "8069", "TCP port to listen on.")

	return flags.Parse(args)
}

func (s *StartCmd) Run() error {
	return http.ListenAndServe(":"+s.port, s.server)
}
