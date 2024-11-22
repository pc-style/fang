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
	codeBlockStyle = lipgloss.NewStyle().Background(bg).MarginLeft(2).Padding(1, 2)
	programStyle   = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("#7E65FF")).PaddingLeft(1)
	commentStyle   = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.BrightBlack).PaddingLeft(1)
	stringStyle    = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("#02BF87")).PaddingLeft(1)
	argumentStyle  = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("248")).PaddingLeft(1)
	flagStyle      = lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("244")).PaddingLeft(1)
	titleStyle     = lipgloss.NewStyle().Bold(true).Transform(strings.ToUpper).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#6C50FF"))

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

	cmds := map[string]string{}
	for _, sc := range c.Commands() {
		cmds[lipgloss.NewStyle().PaddingLeft(3).Render(sc.Use)] = helpStyle.Render(sc.Short)
	}

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
		usage := helpStyle.Render(f.Usage)
		if f.DefValue != "" {
			usage = lipgloss.JoinHorizontal(
				lipgloss.Left,
				usage,
				helpStyle.Render(" ("),
				keywordStyle.Render(f.DefValue),
				helpStyle.Render(")"),
			)
		}
		flags[key] = usage
	})

	space := minSpace
	for _, k := range append(
		slices.Collect(maps.Keys(flags)),
		slices.Collect(maps.Keys(cmds))...,
	) {
		space = max(space, lipgloss.Width(k)+2)
	}

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

	examples := strings.Split(c.Example, "\n")
	for i, line := range examples {
		line = strings.TrimSpace(line)
		if line == "" && i < len(examples)-1 {
			usage = append(usage, "")
			continue
		}

		if strings.HasPrefix(line, "# ") {
			usage = append(usage, commentStyle.Render(line))
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

		usage = append(usage, lipgloss.JoinHorizontal(
			lipgloss.Left,
			args...,
		))
	}

	return usage
}
