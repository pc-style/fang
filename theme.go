package serpentine

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

type Theme struct {
	Codeblock color.Color
	Title     color.Color
	Comment   color.Color
	Argument  color.Color
	Flag      color.Color
	Help      color.Color
	Default   color.Color
	Dash      color.Color
}

var DefaultTheme = Theme{
	Codeblock: lipgloss.Color("235"),
	Title:     lipgloss.Color("141"),
	Comment:   lipgloss.Color("8"),
	Flag:      lipgloss.Color("250"),
	Argument:  lipgloss.Color("248"),
	Help:      lipgloss.Color("243"),
	Dash:      lipgloss.Color("240"),
	Default:   lipgloss.Color("146"),
}

type Styles struct {
	Codeblock lipgloss.Style
	Program   lipgloss.Style
	Comment   lipgloss.Style
	Argument  lipgloss.Style
	Flag      lipgloss.Style
	Title     lipgloss.Style
	Span      lipgloss.Style
	Dash      lipgloss.Style
	Help      lipgloss.Style
	Default   lipgloss.Style
}

func makeStyles(theme Theme) Styles {
	return Styles{
		Codeblock: lipgloss.NewStyle().
			Background(theme.Codeblock).
			MarginLeft(2).
			Padding(1, 3, 0, 1),
		Program: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Title).
			PaddingLeft(1),
		Comment: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Comment).
			PaddingLeft(1),
		Argument: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Argument).
			PaddingLeft(1),
		Flag: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Flag).
			PaddingLeft(1),
		Dash: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Dash).
			PaddingLeft(1),
		Span: lipgloss.NewStyle().
			Background(theme.Codeblock),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Title).
			Transform(strings.ToUpper).
			Margin(1, 0, 0, 2),
		Help: lipgloss.NewStyle().
			Foreground(theme.Help),
		Default: lipgloss.NewStyle().
			Foreground(theme.Default),
	}
}

func (s Styles) nobg() Styles {
	return Styles{
		Codeblock: s.Codeblock.UnsetBackground(),
		Program:   s.Program.UnsetBackground(),
		Comment:   s.Comment.UnsetBackground(),
		Argument:  s.Argument.UnsetBackground(),
		Flag:      s.Flag.UnsetBackground(),
		Dash:      s.Dash.UnsetBackground(),
		Span:      s.Span.UnsetBackground(),
		Help:      s.Help,
	}
}
