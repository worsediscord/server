package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/worsediscord/server/cmd"
)

type RootCmd struct {
	name        string
	subcommands []cmd.Command
}

func NewRootCmd(name string, subcommands ...cmd.Command) *RootCmd {
	if len(name) == 0 {
		if len(os.Args) >= 1 {
			name = filepath.Base(os.Args[0])
		} else {
			name = "root"
		}
	}

	rootCmd := &RootCmd{
		name:        name,
		subcommands: subcommands,
	}

	return rootCmd
}

func (r *RootCmd) Name() string {
	return r.name
}

func (r *RootCmd) Description() string {
	return "Manage a worsediscord server"

}

// Parse parses the global flags of the program. The []string parameter is ignored.
func (r *RootCmd) Parse([]string) error {
	flag.Usage = func() {
		_, _ = fmt.Fprint(flag.CommandLine.Output(), cmd.HelpString("", r, flag.CommandLine, r.subcommands...))
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

func (r *RootCmd) AddSubcommands(subcommands ...cmd.Command) {
	r.subcommands = append(r.subcommands, subcommands...)
}
