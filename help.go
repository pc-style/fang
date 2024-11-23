package serpentine

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const minSpace = 10

// TODO: use ansi colors only
var (
	bg             = lipgloss.Color("235")
	codeBlockStyle = lipgloss.NewStyle().Background(bg).MarginLeft(2).Padding(1, 3, 0, 1)
	programStyle   = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("#7E65FF")).PaddingLeft(1)
	commentStyle   = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.BrightBlack).PaddingLeft(1)
	argumentStyle  = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("248")).PaddingLeft(1)
	flagStyle      = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("244")).PaddingLeft(1)
	titleStyle     = lipgloss.NewStyle().Bold(true).Transform(strings.ToUpper).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#6C50FF"))
	spanStyle      = lipgloss.NewStyle().Background(bg)

	dashStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginLeft(1)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

func helpFn(c *cobra.Command, _ []string) {
	w := colorprofile.NewWriter(os.Stdout, os.Environ())
	_, _ = fmt.Fprintln(w)
	if c.Long == "" {
		_, _ = fmt.Fprintln(w, "  "+c.Short)
	} else {
		_, _ = fmt.Fprintln(w, lipgloss.NewStyle().PaddingLeft(2).Render(c.Long))
	}
	_, _ = fmt.Fprintln(w, titleStyle.Render("usage"))
	_, _ = fmt.Fprintln(w)

	_, _ = fmt.Fprintln(
		w,
		codeBlockStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				usage(c)...,
			),
		),
	)

	cmds, cmdKeys := evalCmds(c)
	flags, flagKeys := evalFlags(c)
	space := calculateSpace(cmdKeys, flagKeys)

	if len(cmds) > 0 {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, titleStyle.Render("commands"))
		for _, k := range cmdKeys {
			_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				cmds[k],
			))
		}
	}

	if len(flags) > 0 {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, titleStyle.Render("flags"))
		for _, k := range flagKeys {
			_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				flags[k],
			))
		}
	}
}

var otherArgsRe = regexp.MustCompile(`(\[.*\])`)

// use stylized use line for a given command.
func use(c *cobra.Command, inline bool) string {
	u := c.Use
	hasArgs := strings.Contains(u, "[args]")
	hasFlags := strings.Contains(u, "[flags]") || strings.Contains(u, "[--flags]") || c.HasFlags() || c.HasPersistentFlags() || c.HasAvailableFlags()
	hasCommands := strings.Contains(u, "[command]") || c.HasAvailableSubCommands()
	for _, k := range []string{
		"[args]",
		"[flags]", "[--flags]",
		"[command]",
	} {
		u = strings.ReplaceAll(u, k, "")
	}

	var otherArgs []string
	for _, arg := range otherArgsRe.FindAllString(u, -1) {
		u = strings.ReplaceAll(u, arg, "")
		otherArgs = append(otherArgs, arg)
	}

	u = strings.TrimSpace(u)

	programStyle := programStyle
	argumentStyle := argumentStyle
	flagStyle := flagStyle
	if inline {
		programStyle = programStyle.UnsetBackground()
		argumentStyle = argumentStyle.UnsetBackground()
		flagStyle = flagStyle.UnsetBackground()
	}
	useLine := []string{programStyle.Render(u)}
	if hasCommands {
		useLine = append(useLine, argumentStyle.Render("[command]"))
	}
	if hasArgs {
		useLine = append(useLine, argumentStyle.Render("[args]"))
	}
	for _, arg := range otherArgs {
		useLine = append(useLine, argumentStyle.Render(arg))
	}
	if hasFlags {
		useLine = append(useLine, flagStyle.Render("[--flags]"))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, useLine...)
}

// usage for a given command.
// will print both the cmd.Use and cmd.Example bits.
func usage(c *cobra.Command) []string {
	usage := []string{use(c, false)}

	size := lipgloss.Width(usage[0])
	examples := strings.Split(c.Example, "\n")
	for i, line := range examples {
		line = strings.TrimSpace(line)
		if line == "" && i != len(examples)-1 {
			usage = append(usage, " ")
			continue
		}

		if strings.HasPrefix(line, "# ") {
			s := lipgloss.JoinHorizontal(
				lipgloss.Left,
				commentStyle.Render(line),
			)
			size = max(size, lipgloss.Width(s))
			usage = append(usage, s)
			continue
		}

		args := strings.Fields(line)
		for i, arg := range args {
			if i == 0 {
				args[i] = programStyle.Render(arg)
				continue
			}
			args[i] = argumentStyle.Render(arg)
		}

		s := lipgloss.JoinHorizontal(
			lipgloss.Left,
			args...,
		)
		size = max(size, lipgloss.Width(s))
		usage = append(usage, s)
	}

	for i, use := range usage {
		if usize := lipgloss.Width(use); size > usize {
			usage[i] = lipgloss.JoinHorizontal(
				lipgloss.Left,
				use,
				spanStyle.Render(strings.Repeat(" ", size-usize)),
			)
		}
	}
	return usage
}

func evalFlags(c *cobra.Command) (map[string]string, []string) {
	flags := map[string]string{}
	keys := []string{}
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		var parts []string
		if f.Shorthand == "" {
			parts = append(
				parts,
				dashStyle.PaddingLeft(5).Render("--"),
				f.Name,
			)
		} else {
			parts = append(
				parts,
				dashStyle.PaddingLeft(2).Render("-"),
				f.Shorthand,
				dashStyle.Render("--"),
				f.Name,
			)
		}
		key := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
		help := helpStyle.Render(f.Usage)
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
			help = lipgloss.JoinHorizontal(
				lipgloss.Left,
				helpStyle.Render(f.Usage+" ("),
				keywordStyle.Render(f.DefValue),
				helpStyle.Render(")"),
			)
		}
		flags[key] = help
		keys = append(keys, key)
	})
	return flags, keys
}

func evalCmds(c *cobra.Command) (map[string]string, []string) {
	keys := []string{}
	pad := lipgloss.NewStyle().PaddingLeft(3)
	cmds := map[string]string{}
	for _, sc := range c.Commands() {
		if sc.Hidden {
			continue
		}
		key := pad.Render(use(sc, true))
		help := helpStyle.Render(sc.Short)
		cmds[key] = help
		keys = append(keys, key)
	}
	return cmds, keys
}

func calculateSpace(k1, k2 []string) int {
	space := minSpace
	for _, k := range append(k1, k2...) {
		space = max(space, lipgloss.Width(k)+2)
	}
	return space
}
