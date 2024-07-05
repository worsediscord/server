package cmd

import (
	"flag"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"
)

const helpLongFlagTag = "[long]"
const helpFlagDelim = ":"
const padding = "  "

type helpFlag struct {
	shortFlag string
	longFlag  string
	usage     string
}

// LongFlagUsage returns a formatted string to be used in flag.Flag's usage string as a long version.
func LongFlagUsage(shortFlagName string) string {
	return fmt.Sprintf("%s%s%s", helpLongFlagTag, helpFlagDelim, shortFlagName)
}

func HelpString(commandPrefix string, command Command, fs *flag.FlagSet, subcommands ...Command) string {
	var helpString strings.Builder
	w := tabwriter.NewWriter(&helpString, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintf(w, "%s\n\nUsage:\n%s%s%s", command.Description(), padding, commandPrefix, command.Name())

	if fs != nil && countFlags(fs) > 0 {
		_, _ = fmt.Fprintf(w, " [options]")
	}

	if len(subcommands) > 0 {
		_, _ = fmt.Fprintf(w, " [command]\n\nAvailable Commands\n")

		for _, cmd := range subcommands {
			_, _ = fmt.Fprintf(w, "%s%s\t%s\n", padding, cmd.Name(), cmd.Description())
		}
	} else {
		_, _ = fmt.Fprint(w, "\n")
	}

	if fs == nil {
		_ = w.Flush()
		return helpString.String()
	}

	var helpFlags []helpFlag
	var orphanedFlags []*flag.Flag
	fs.VisitAll(func(f *flag.Flag) {
		if strings.HasPrefix(f.Usage, helpLongFlagTag) {
			shortFlag := getShortFlag(f.Usage)
			n := slices.IndexFunc(helpFlags, func(h helpFlag) bool {
				return h.shortFlag == shortFlag
			})

			if n != -1 && len(shortFlag) > 0 {
				helpFlags[n].longFlag = f.Name
			} else {
				orphanedFlags = append(orphanedFlags, f)
			}
		} else {
			if len(f.Name) == 1 {
				helpFlags = append(helpFlags, helpFlag{shortFlag: f.Name, usage: f.Usage})
			} else {
				helpFlags = append(helpFlags, helpFlag{longFlag: f.Name, usage: f.Usage})
			}
		}
	})

	// If for some reason fs wasn't nil, but we didn't find any formal flags, return early
	if len(helpFlags) == 0 {
		_ = w.Flush()
		return helpString.String()
	}

	// Loop through any orphaned flags once more to see if we just got unlucky with ordering
	for _, f := range orphanedFlags {
		if strings.HasPrefix(f.Usage, helpLongFlagTag) {
			shortFlag := getShortFlag(f.Usage)
			n := slices.IndexFunc(helpFlags, func(h helpFlag) bool {
				return h.shortFlag == shortFlag
			})
			if n != -1 {
				helpFlags[n].longFlag = f.Name
			}
		}
	}

	_, _ = fmt.Fprintf(w, "\nOptions:\n")
	for _, h := range helpFlags {
		_, _ = fmt.Fprint(w, padding)

		if h.shortFlag != "" {
			_, _ = fmt.Fprintf(w, "-%s", h.shortFlag)
		} else {
			// pad four spaces for: - + [short flag] + , + ' '
			_, _ = fmt.Fprint(w, "    ")
		}

		if h.longFlag != "" {
			if h.shortFlag != "" {
				_, _ = fmt.Fprint(w, ", ")
			}

			_, _ = fmt.Fprintf(w, "--%s", h.longFlag)
		}

		_, _ = fmt.Fprintf(w, "\t\t%s\n", h.usage)
	}

	_ = w.Flush()

	return helpString.String()
}

func getShortFlag(longFlagUsage string) string {
	split := strings.Split(longFlagUsage, helpFlagDelim)

	if len(split) != 2 {
		return ""
	}

	return split[1]
}

// countFlag counts all flags in a flag set, even those not set.
func countFlags(fs *flag.FlagSet) int {
	count := 0
	fs.VisitAll(func(f *flag.Flag) {
		count++
	})

	return count
}
