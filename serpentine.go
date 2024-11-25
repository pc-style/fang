package serpentine

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/colorprofile"
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
	theme       Theme
}

// Option changes serpentine settings.
type Option func(*settings)

// WithoutCompletions disables completions.
func WithoutCompletions() Option {
	return func(s *settings) {
		s.completions = false
	}
}

// WithoutManpage disables man pages.
func WithoutManpage() Option {
	return func(s *settings) {
		s.manpages = false
	}
}

// WithTheme sets the colorscheme.
func WithTheme(theme Theme) Option {
	return func(s *settings) {
		s.theme = theme
	}
}

// WithVersion sets the version.
func WithVersion(version string) Option {
	return func(s *settings) {
		s.version = version
	}
}

// WithCommit sets the commit SHA.
func WithCommit(commit string) Option {
	return func(s *settings) {
		s.commit = commit
	}
}

type Command interface {
	Execute() error
}

type cobraCmd struct {
	*cobra.Command
	styles Styles
}

func (c *cobraCmd) Execute() error {
	if err := c.Command.Execute(); err != nil {
		w := colorprofile.NewWriter(c.ErrOrStderr(), os.Environ())
		writeError(w, c.styles, err)
		return err
	}
	return nil
}

// Setup setups the given root *cobra.Command.
func Setup(root *cobra.Command, options ...Option) Command {
	opts := settings{
		manpages:    true,
		completions: true,
		theme:       DefaultTheme,
	}
	for _, option := range options {
		option(&opts)
	}

	styles := makeStyles(opts.theme)

	root.SetHelpFunc(func(c *cobra.Command, _ []string) {
		w := colorprofile.NewWriter(c.OutOrStdout(), os.Environ())
		helpFn(c, w, styles)
	})
	root.SilenceUsage = true
	root.SilenceErrors = true

	if opts.manpages {
		root.AddCommand(&cobra.Command{
			Use:                   "man",
			Short:                 "Generates manpages",
			SilenceUsage:          true,
			DisableFlagsInUseLine: true,
			Hidden:                true,
			Args:                  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				page, err := mango.NewManPage(1, cmd.Root())
				if err != nil {
					//nolint:wrapcheck
					return err
				}
				_, err = fmt.Fprint(os.Stdout, page.Build(roff.NewDocument()))
				//nolint:wrapcheck
				return err
			},
		})
	}

	if opts.completions {
		root.InitDefaultCompletionCmd()
	} else {
		root.CompletionOptions.DisableDefaultCmd = true
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

	root.Version = opts.version

	return &cobraCmd{
		Command: root,
		styles:  styles,
	}
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
