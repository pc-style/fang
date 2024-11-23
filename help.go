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

func helpFn(c *cobra.Command, styles Styles) {
	const shortPad = 2
	w := colorprofile.NewWriter(os.Stdout, os.Environ())
	_, _ = fmt.Fprintln(w)
	if c.Long == "" {
		_, _ = fmt.Fprintln(w, "  "+c.Short)
	} else {
		_, _ = fmt.Fprintln(w, lipgloss.NewStyle().PaddingLeft(shortPad).Render(c.Long))
	}
	_, _ = fmt.Fprintln(w, styles.Title.Render("usage"))
	_, _ = fmt.Fprintln(w)

	_, _ = fmt.Fprintln(
		w,
		styles.Codeblock.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				usage(c, styles)...,
			),
		),
	)

	cmds, cmdKeys := evalCmds(c, styles)
	flags, flagKeys := evalFlags(c, styles)
	space := calculateSpace(cmdKeys, flagKeys)

	if len(cmds) > 0 {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, styles.Title.Render("commands"))
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
		_, _ = fmt.Fprintln(w, styles.Title.Render("flags"))
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
//
//nolint:mnd
func use(c *cobra.Command, styles Styles, inline bool) string {
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

	programStyle := styles.Program
	argumentStyle := styles.Argument
	flagStyle := styles.Flag
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
//
//nolint:mnd
func usage(c *cobra.Command, styles Styles) []string {
	usage := []string{use(c, styles, false)}

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
				styles.Comment.Render(line),
			)
			size = max(size, lipgloss.Width(s))
			usage = append(usage, s)
			continue
		}

		args := strings.Fields(line)
		for i, arg := range args {
			if i == 0 {
				args[i] = styles.Program.Render(arg)
				continue
			}
			args[i] = styles.Argument.Render(arg)
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
				styles.Span.Render(strings.Repeat(" ", size-usize)),
			)
		}
	}
	return usage
}

func evalFlags(c *cobra.Command, styles Styles) (map[string]string, []string) {
	const shortPad = 2
	const noShortPad = shortPad + 3
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
				styles.Dash.PaddingLeft(noShortPad).Render("--"),
				f.Name,
			)
		} else {
			parts = append(
				parts,
				styles.Dash.PaddingLeft(shortPad).Render("-"),
				f.Shorthand,
				styles.Dash.Render("--"),
				f.Name,
			)
		}
		key := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
		help := styles.Help.Render(f.Usage)
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
			help = lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.Help.Render(f.Usage+" ("),
				styles.Default.Render(f.DefValue),
				styles.Help.Render(")"),
			)
		}
		flags[key] = help
		keys = append(keys, key)
	})
	return flags, keys
}

func evalCmds(c *cobra.Command, styles Styles) (map[string]string, []string) {
	padStyle := lipgloss.NewStyle().PaddingLeft(3) //nolint:mnd
	keys := []string{}
	cmds := map[string]string{}
	for _, sc := range c.Commands() {
		if sc.Hidden {
			continue
		}
		key := padStyle.Render(use(sc, styles, true))
		help := styles.Help.Render(sc.Short)
		cmds[key] = help
		keys = append(keys, key)
	}
	return cmds, keys
}

func calculateSpace(k1, k2 []string) int {
	const spaceBetween = 2
	space := minSpace
	for _, k := range append(k1, k2...) {
		space = max(space, lipgloss.Width(k)+spaceBetween)
	}
	return space
}
