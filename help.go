package serpentine

import (
	"fmt"
	"maps"
	"os"
	"slices"
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
	stringStyle    = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("#02BF87")).PaddingLeft(1)
	argumentStyle  = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("248")).PaddingLeft(1)
	flagStyle      = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("244")).PaddingLeft(1)
	titleStyle     = lipgloss.NewStyle().Bold(true).Transform(strings.ToUpper).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#6C50FF"))
	spanStyle      = lipgloss.NewStyle().Background(bg)

	dashStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginLeft(1)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

func helpFn(c *cobra.Command, args []string) {
	w := colorprofile.NewWriter(os.Stdout, os.Environ())
	fmt.Fprintln(w)
	if c.Long == "" {
		fmt.Fprintln(w, "  "+c.Short)
	} else {
		fmt.Fprintln(w, lipgloss.NewStyle().PaddingLeft(2).Render(c.Long))
	}
	fmt.Fprintln(w, titleStyle.Render("usage"))
	fmt.Fprintln(w)

	fmt.Fprintln(
		w,
		codeBlockStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				usage(c)...,
			),
		),
	)

	cmds, flags := evalCmds(c), evalFlags(c)
	space := calculateSpace(cmds, flags)

	if len(cmds) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, titleStyle.Render("commands"))
		for k, v := range cmds {
			fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				v,
			))
		}
	}

	if len(flags) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, titleStyle.Render("flags"))
		for k, v := range flags {
			fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				v,
			))
		}
	}
}

func usage(c *cobra.Command) []string {
	useLine := []string{programStyle.Render(c.Use)}
	if c.HasSubCommands() {
		useLine = append(useLine, argumentStyle.Render("[command]"))
	}
	if c.HasFlags() || c.HasPersistentFlags() {
		useLine = append(useLine, flagStyle.Render("[--flags]"))
	}

	usage := []string{
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			useLine...,
		),
	}

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

func evalFlags(c *cobra.Command) map[string]string {
	flags := map[string]string{}
	c.Flags().VisitAll(func(f *pflag.Flag) {
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
		if f.DefValue != "" {
			help = lipgloss.JoinHorizontal(
				lipgloss.Left,
				help,
				helpStyle.Render(" ("),
				keywordStyle.Render(f.DefValue),
				helpStyle.Render(")"),
			)
		}
		flags[key] = help
	})
	return flags
}

func evalCmds(c *cobra.Command) map[string]string {
	pad := lipgloss.NewStyle().PaddingLeft(3)
	cmds := map[string]string{}
	for _, sc := range c.Commands() {
		key := pad.Render(sc.Use)
		// handles native commands, such as 'help', which report use as `help [command]`.
		if strings.Contains(key, "[command]") {
			key = lipgloss.JoinHorizontal(
				lipgloss.Left,
				pad.Render(strings.TrimSuffix(sc.Use, " [command]")),
				argumentStyle.UnsetBackground().Render("[command]"),
			)
		}
		help := helpStyle.Render(sc.Short)
		cmds[key] = help
	}
	return cmds
}

func calculateSpace(m1, m2 map[string]string) int {
	space := minSpace
	for _, k := range append(
		slices.Collect(maps.Keys(m1)),
		slices.Collect(maps.Keys(m2))...,
	) {
		space = max(space, lipgloss.Width(k)+2)
	}
	return space
}
