package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/pcstyle/kawaii-shell/internal/shell"
	"github.com/pcstyle/kawaii-shell/internal/themes"
	"github.com/pcstyle/kawaii-shell/internal/ui/components"
	"github.com/pcstyle/kawaii-shell/internal/ui/pet"
)

// App is the main Bubble Tea application model
type App struct {
	shell       *shell.Shell
	pet         *pet.Pet
	theme       *themes.KawaiiTheme
	startup     *components.StartupSequence
	input       string
	output      []string
	prompt      string
	cursor      int
	width       int
	height      int
	ready       bool
	lastCommand string
}

// NewApp creates a new kawaii shell application
func NewApp() *App {
	sh, _ := shell.NewShell()

	return &App{
		shell:  sh,
		pet:    pet.NewPet("Neko", pet.TypeCat),
		theme:  themes.NewSakuraTheme(),
		prompt: "🌸> ",
		output: []string{
			"✨ Welcome to Kawaii Shell! ✨",
			"Your adorable terminal companion! 🐱",
			"",
			"Type 'help' for cute commands, or any regular command!",
		},
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Start the shell
	if err := a.shell.Start(); err != nil {
		a.output = append(a.output, "🥺 Oops! Couldn't start shell: "+err.Error())
	}

	return tea.Batch(
		a.pet.Init(),
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return TickMsg{Time: t}
		}),
	)
}

// TickMsg represents a tick for animations
type TickMsg struct {
	Time time.Time
}

// Update handles messages and updates the model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		if a.startup == nil {
			a.startup = components.NewStartupSequence(a.width, a.height, "0.1.0")
			cmds = append(cmds, a.startup.Init())
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit

		case "enter":
			if strings.TrimSpace(a.input) != "" {
				a.executeCommand(a.input)
				a.lastCommand = a.input
				a.input = ""
				a.cursor = 0
			}

		case "backspace":
			if a.cursor > 0 {
				a.input = a.input[:a.cursor-1] + a.input[a.cursor:]
				a.cursor--
			}

		case "left":
			if a.cursor > 0 {
				a.cursor--
			}

		case "right":
			if a.cursor < len(a.input) {
				a.cursor++
			}

		default:
			if len(msg.String()) == 1 {
				a.input = a.input[:a.cursor] + msg.String() + a.input[a.cursor:]
				a.cursor++
			}
		}

	case TickMsg:
		var petCmd tea.Cmd
		a.pet, petCmd = a.pet.Update(msg)
		if petCmd != nil {
			cmds = append(cmds, petCmd)
		}
		if output := a.shell.GetOutput(); output != "" {
			a.output = append(a.output, output)
		}
		cmds = append(cmds, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return TickMsg{Time: t}
		}))
	case components.StartupTickMsg:
		if a.startup != nil && !a.startup.IsComplete() {
			var cmd tea.Cmd
			a.startup, cmd = a.startup.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case components.ParticleTickMsg:
		if a.startup != nil && !a.startup.IsComplete() {
			var cmd tea.Cmd
			a.startup, cmd = a.startup.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return a, tea.Batch(cmds...)
}

// executeCommand processes and executes a command
func (a *App) executeCommand(command string) {
	// Get cute command info
	info := shell.GetCommandInfo(command)

	// Show what we're doing
	cuteLine := a.theme.Styles.CommandInfo.Render(
		info.Emoji + " " + info.FriendlyName + ": " + info.Description,
	)
	a.output = append(a.output, cuteLine)

	// Show warning for dangerous commands
	if info.IsDangerous {
		warning := a.theme.Styles.Warning.Render(
			"⚠️  This command might be dangerous! Are you sure? (This is just a demo - command will execute)",
		)
		a.output = append(a.output, warning)
	}

	// Update pet reaction
	a.pet.ReactToCommand(command, info.IsDangerous)

	// Handle special kawaii commands
	switch strings.TrimSpace(command) {
	case "help":
		a.showHelp()
		return
	case "kawaii":
		a.showKawaii()
		return
	case "pet":
		a.showPetStatus()
		return
	}

	// Execute the actual command
	if err := a.shell.ExecuteCommand(command); err != nil {
		a.output = append(a.output, "🥺 Oops: "+err.Error())
	}
}

// showHelp displays cute help information
func (a *App) showHelp() {
	help := []string{
		"",
		"🌸 ✨ Kawaii Shell Commands ✨ 🌸",
		"",
		"🐱 kawaii    - Show kawaii info",
		"🐱 pet       - Check your pet's status",
		"🐱 help      - Show this cute help",
		"",
		"✨ All regular commands work too! ✨",
		"I'll make them cute and friendly! 💕",
		"",
	}

	for _, line := range help {
		a.output = append(a.output, a.theme.Styles.Help.Render(line))
	}
}

// showKawaii displays kawaii information
func (a *App) showKawaii() {
	kawaii := []string{
		"",
		"🌸 ✨ About Kawaii Shell ✨ 🌸",
		"",
		"Making terminals cute and friendly! 💕",
		"Your commands become adorable adventures!",
		"",
		fmt.Sprintf("🐱 Pet: %s (%s)", a.pet.Name, a.pet.GetMoodEmoji()),
		fmt.Sprintf("🎨 Theme: %s", a.theme.Name),
		"",
		"Built with lots of love and sparkles! ✨",
		"",
	}

	for _, line := range kawaii {
		a.output = append(a.output, a.theme.Styles.Info.Render(line))
	}
}

// showPetStatus shows the pet's current status
func (a *App) showPetStatus() {
	status := a.pet.GetStatus()
	for _, line := range status {
		a.output = append(a.output, a.theme.Styles.Pet.Render(line))
	}
}

// View renders the application
func (a *App) View() string {
	if !a.ready {
		return "Loading kawaii shell... ✨"
	}
	if a.startup != nil && !a.startup.IsComplete() {
		return a.startup.Render()
	}
	petHeight := 4
	inputHeight := 3
	availableHeight := a.height - petHeight - inputHeight - 2
	outputLines := a.output
	if len(outputLines) > availableHeight {
		outputLines = outputLines[len(outputLines)-availableHeight:]
	}
	output := strings.Join(outputLines, "\n")
	outputBox := a.theme.Styles.OutputBox.
		Width(a.width - 2).
		Height(availableHeight).
		Render(output)
	inputText := a.input
	if a.cursor < len(inputText) {
		inputText = inputText[:a.cursor] + a.theme.Styles.Cursor.Render(string(inputText[a.cursor])) + inputText[a.cursor+1:]
	} else {
		inputText += a.theme.Styles.Cursor.Render(" ")
	}
	inputLine := a.theme.Styles.Prompt.Render(a.prompt) + a.theme.Styles.Input.Render(inputText)
	inputBox := a.theme.Styles.InputBox.
		Width(a.width - 2).
		Render(inputLine)
	petView := a.pet.View()
	petBox := a.theme.Styles.PetBox.
		Width(20).
		Height(petHeight).
		Render(petView)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		outputBox,
		"",
		inputBox,
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainContent,
		lipgloss.NewStyle().Width(a.width-len(mainContent)).Render(""),
	) + "\n" + lipgloss.PlaceHorizontal(a.width, lipgloss.Right, petBox)
}
