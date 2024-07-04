package main

import (
	"flag"
	"fmt"
	"github.com/worsediscord/server/cmd"
	"os"

	"github.com/worsediscord/server/api"
)

type RootCmd struct {
	server      *api.Server
	subcommands []cmd.Command

	verbose bool
}

func NewRootCmd(subcommands ...cmd.Command) *RootCmd {
	rootCmd := &RootCmd{
		subcommands: subcommands,
	}

	return rootCmd
}

func (r *RootCmd) Name() string {
	var name string
	if len(os.Args) >= 1 {
		name = os.Args[0]
	}

	return name
}

func (r *RootCmd) Description() string {
	return "Manage a worsediscord server"

}

// Parse parses the global flags of the program. The []string parameter is ignored.
func (r *RootCmd) Parse([]string) error {
	flag.Usage = func() {
		_, _ = fmt.Fprint(flag.CommandLine.Output(), cmd.HelpString(r, flag.CommandLine, r.subcommands...))
	}

	flag.Parse()

	return nil
}

func (r *RootCmd) Run() error {
	for _, subcommand := range r.subcommands {
		if flag.Arg(0) == subcommand.Name() {
			if err := subcommand.Parse(flag.Args()[1:]); err != nil {
				return err
			}

			if err := subcommand.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
