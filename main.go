package serpentine

import (
	"fmt"
	"os"
	"runtime/debug"

	mango "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

const shaLen = 7

type settings struct {
	completions bool
	manpages    bool
	version     string
	commit      string
}

type Opt func(*settings)

func WithoutCompletions(s *settings) {
	s.completions = false
}

func WithoutManPage(s *settings) {
	s.manpages = false
}

func WithVersion(version, commit string) Opt {
	return func(s *settings) {
		s.version = version
		s.commit = commit
	}
}

func Setup(cmd *cobra.Command, options ...Opt) *cobra.Command {
	opts := settings{
		manpages:    true,
		completions: true,
	}
	for _, option := range options {
		option(&opts)
	}

	cmd.SetHelpFunc(helpFn)

	if opts.manpages {
		cmd.AddCommand(&cobra.Command{
			Use:                   "man",
			Short:                 "Generates manpages",
			SilenceUsage:          true,
			DisableFlagsInUseLine: true,
			Hidden:                true,
			Args:                  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				page, err := mango.NewManPage(1, cmd.Root())
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(os.Stdout, page.Build(roff.NewDocument()))
				return err
			},
		})
	}

	if opts.completions {
		cmd.InitDefaultCompletionCmd()
	}

	if opts.version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			opts.version = info.Main.Version
			opts.commit = getKey(info, "vcs.revision")
		} else {
			opts.version = "unknown (built from source)"
		}
	}
	if len(opts.commit) >= shaLen {
		opts.version += " (" + opts.commit[:shaLen] + ")"
	}

	cmd.Version = opts.version

	return cmd
}

func getKey(info *debug.BuildInfo, key string) string {
	if info == nil {
		return ""
	}
	for _, iter := range info.Settings {
		if iter.Key == key {
			return iter.Value
		}
	}
	return ""
}
