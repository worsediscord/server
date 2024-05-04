package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/worsediscord/server/api"
)

type RootCmd struct {
	server      *api.Server
	subcommands []Command

	verbose bool
}

func NewRootCmd(subcommands ...Command) *RootCmd {

	rootCmd := &RootCmd{
		subcommands: subcommands,
	}

	helpString := fmt.Sprintf("%s\n\nUsage:\n  %s [options] [command]\n", rootCmd.Description(), rootCmd.Name())
	if len(rootCmd.subcommands) > 0 {
		helpString += "\nAvailable Commands:\n"

		for _, cmd := range rootCmd.subcommands {
			helpString += fmt.Sprintf("  %s    \t%s\n", cmd.Name(), cmd.Description())
		}
	}
	helpString += fmt.Sprintf("\nOptions:\n")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), helpString)
		flag.PrintDefaults()
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
	flag.BoolVar(&r.verbose, "v", false, "Enable verbose logging.")
	flag.Parse()

	return nil
}

func (r *RootCmd) Run() error {
	for _, cmd := range r.subcommands {
		if flag.Arg(0) == cmd.Name() {
			if err := cmd.Parse(flag.Args()[1:]); err != nil {
				return err
			}

			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
