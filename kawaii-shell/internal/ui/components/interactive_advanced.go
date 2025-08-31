package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// Tab represents a single tab in a tab group
type Tab struct {
	Title     string
	Content   string
	Active    bool
	Particles *ParticleSystem
}

// TabGroup represents a group of interactive tabs
type TabGroup struct {
	X, Y        int
	Width       int
	Height      int
	Tabs        []Tab
	ActiveTab   int
	Style       lipgloss.Style
	ActiveStyle lipgloss.Style
	Focused     bool
	Animation   *AnimatedElement
}

// NewTabGroup creates a stunning tab group
func NewTabGroup(x, y, width, height int) *TabGroup {
	baseStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Guppy).
		Background(lipgloss.Color("#f0f8ff")).
		Padding(1)

	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(charmtone.Coral).
		Background(lipgloss.Color("#fff8f8")).
		Foreground(charmtone.Charcoal).
		Bold(true).
		Transform(func(s string) string {
			return fmt.Sprintf("✨ %s ✨", s)
		})

	return &TabGroup{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Tabs:        make([]Tab, 0),
		ActiveTab:   0,
		Style:       baseStyle,
		ActiveStyle: activeStyle,
		Animation:   NewAnimatedElement("", float64(x), float64(y)),
	}
}

// AddTab adds a new tab
func (tg *TabGroup) AddTab(title, content string) {
	tab := Tab{
		Title:     title,
		Content:   content,
		Active:    len(tg.Tabs) == 0,
		Particles: NewParticleSystem(20, 10),
	}
	tg.Tabs = append(tg.Tabs, tab)
}

// Update updates the tab group
func (tg *TabGroup) Update(msg tea.Msg) (*TabGroup, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		for i := range tg.Tabs {
			tg.Tabs[i].Particles.Update(0.05)
		}
		cmds = append(cmds, ParticleUpdateCmd())

	case tea.KeyMsg:
		if tg.Focused {
			switch msg.String() {
			case "tab", "right", "l":
				tg.NextTab()
				cmds = append(cmds, tg.createTabSwitchEffect())
			case "shift+tab", "left", "h":
				tg.PrevTab()
				cmds = append(cmds, tg.createTabSwitchEffect())
			}
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Check if clicking on tab headers
			for i, tab := range tg.Tabs {
				tabX := tg.X + i*15
				if msg.X >= tabX && msg.X <= tabX+len(tab.Title)+4 && msg.Y == tg.Y {
					tg.SetActiveTab(i)
					cmds = append(cmds, tg.createTabSwitchEffect())
					break
				}
			}
		}
	}

	return tg, tea.Batch(cmds...)
}

// NextTab switches to next tab
func (tg *TabGroup) NextTab() {
	if len(tg.Tabs) == 0 {
		return
	}
	tg.SetActiveTab((tg.ActiveTab + 1) % len(tg.Tabs))
}

// PrevTab switches to previous tab
func (tg *TabGroup) PrevTab() {
	if len(tg.Tabs) == 0 {
		return
	}
	tg.SetActiveTab((tg.ActiveTab - 1 + len(tg.Tabs)) % len(tg.Tabs))
}

// SetActiveTab sets the active tab
func (tg *TabGroup) SetActiveTab(index int) {
	if index < 0 || index >= len(tg.Tabs) {
		return
	}

	// Reset all tabs
	for i := range tg.Tabs {
		tg.Tabs[i].Active = false
	}

	// Set active tab
	tg.ActiveTab = index
	tg.Tabs[index].Active = true

	// Create celebration effect
	tg.Tabs[index].Particles.AddSparkles(tg.X+index*15+5, tg.Y, 8)
}

// Focus sets tab group focus
func (tg *TabGroup) Focus() {
	tg.Focused = true
}

// Blur removes tab group focus
func (tg *TabGroup) Blur() {
	tg.Focused = false
}

// createTabSwitchEffect creates particle effect when switching tabs
func (tg *TabGroup) createTabSwitchEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		if tg.ActiveTab < len(tg.Tabs) {
			tg.Tabs[tg.ActiveTab].Particles.AddSparkles(tg.X+tg.ActiveTab*15+5, tg.Y, 5)
		}
		return ParticleTickMsg{}
	})
}

// Render renders the stunning tab group
func (tg *TabGroup) Render() string {
	if len(tg.Tabs) == 0 {
		return tg.Style.Width(tg.Width).Height(tg.Height).Render("No tabs")
	}

	// Render tab headers
	var headers []string
	for i, tab := range tg.Tabs {
		style := tg.Style.Copy().
			Padding(0, 2).
			Margin(0, 1).
			Background(lipgloss.Color("#e8f4ff"))

		if tab.Active {
			style = tg.ActiveStyle.Copy().
				Padding(0, 2).
				Margin(0, 1)
		}

		if tg.Focused && tab.Active {
			style = style.Copy().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(charmtone.Malibu)
		}

		headers = append(headers, style.Render(tab.Title))
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)

	// Render active tab content
	var content string
	if tg.ActiveTab < len(tg.Tabs) {
		contentStyle := tg.Style.Copy().
			Width(tg.Width - 2).
			Height(tg.Height - 4).
			Padding(1).
			Background(lipgloss.Color("#ffffff"))

		if tg.Tabs[tg.ActiveTab].Active {
			contentStyle = contentStyle.Copy().
				BorderForeground(charmtone.Coral).
				Background(lipgloss.Color("#fff8f8"))
		}

		content = contentStyle.Render(tg.Tabs[tg.ActiveTab].Content)
	}

	return lipgloss.JoinVertical(lipgloss.Left, headerRow, content)
}

// DropdownOption represents an option in a dropdown menu
type DropdownOption struct {
	Text     string
	Value    interface{}
	Selected bool
}

// Dropdown represents an interactive dropdown menu
type Dropdown struct {
	X, Y          int
	Width         int
	Label         string
	Options       []DropdownOption
	Selected      int
	Open          bool
	Style         lipgloss.Style
	OptionStyle   lipgloss.Style
	SelectedStyle lipgloss.Style
	Focused       bool
	Particles     *ParticleSystem
	Animation     *AnimatedElement
}

// NewDropdown creates a stunning dropdown menu
func NewDropdown(label string, x, y, width int) *Dropdown {
	baseStyle := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Guac).
		Background(lipgloss.Color("#f5fff5")).
		Padding(0, 1)

	optionStyle := lipgloss.NewStyle().
		Width(width-2).
		Padding(0, 1).
		Background(lipgloss.Color("#ffffff"))

	selectedStyle := optionStyle.Copy().
		Background(charmtone.Coral).
		Foreground(charmtone.Butter).
		Bold(true).
		Transform(func(s string) string {
			return fmt.Sprintf("→ %s ←", s)
		})

	return &Dropdown{
		X:             x,
		Y:             y,
		Width:         width,
		Label:         label,
		Options:       make([]DropdownOption, 0),
		Selected:      0,
		Style:         baseStyle,
		OptionStyle:   optionStyle,
		SelectedStyle: selectedStyle,
		Particles:     NewParticleSystem(width+10, 20),
		Animation:     NewAnimatedElement(label, float64(x), float64(y)),
	}
}

// AddOption adds an option to the dropdown
func (d *Dropdown) AddOption(text string, value interface{}) {
	option := DropdownOption{
		Text:     text,
		Value:    value,
		Selected: len(d.Options) == 0,
	}
	d.Options = append(d.Options, option)
}

// Update updates the dropdown
func (d *Dropdown) Update(msg tea.Msg) (*Dropdown, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		d.Particles.Update(0.05)
		cmds = append(cmds, ParticleUpdateCmd())

	case tea.KeyMsg:
		if d.Focused {
			switch msg.String() {
			case "enter", " ":
				d.Toggle()
				cmds = append(cmds, d.createToggleEffect())
			case "up", "k":
				if d.Open && d.Selected > 0 {
					d.Selected--
					d.updateSelection()
					cmds = append(cmds, d.createSelectionEffect())
				}
			case "down", "j":
				if d.Open && d.Selected < len(d.Options)-1 {
					d.Selected++
					d.updateSelection()
					cmds = append(cmds, d.createSelectionEffect())
				}
			case "esc":
				if d.Open {
					d.Close()
				}
			}
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if d.isMouseOverHeader(msg.X, msg.Y) {
				d.Toggle()
				cmds = append(cmds, d.createToggleEffect())
			} else if d.Open {
				optionIndex := d.getOptionFromMouse(msg.X, msg.Y)
				if optionIndex >= 0 {
					d.Selected = optionIndex
					d.updateSelection()
					d.Close()
					cmds = append(cmds, d.createSelectionEffect())
				}
			}
		}
	}

	return d, tea.Batch(cmds...)
}

// Toggle toggles dropdown open/close
func (d *Dropdown) Toggle() {
	d.Open = !d.Open
	if d.Open {
		d.Animation.Bounce(2, 0.5)
	}
}

// Close closes the dropdown
func (d *Dropdown) Close() {
	d.Open = false
}

// updateSelection updates selected option
func (d *Dropdown) updateSelection() {
	for i := range d.Options {
		d.Options[i].Selected = i == d.Selected
	}
}

// isMouseOverHeader checks if mouse is over dropdown header
func (d *Dropdown) isMouseOverHeader(x, y int) bool {
	return x >= d.X && x <= d.X+d.Width && y == d.Y
}

// getOptionFromMouse gets option index from mouse position
func (d *Dropdown) getOptionFromMouse(x, y int) int {
	if x < d.X || x > d.X+d.Width {
		return -1
	}

	optionY := d.Y + 1
	for i := range d.Options {
		if y == optionY+i {
			return i
		}
	}
	return -1
}

// Focus sets dropdown focus
func (d *Dropdown) Focus() {
	d.Focused = true
	d.Animation.Glow(0.5, 2.0)
}

// Blur removes dropdown focus
func (d *Dropdown) Blur() {
	d.Focused = false
	d.Close()
}

// GetSelectedValue returns the selected option value
func (d *Dropdown) GetSelectedValue() interface{} {
	if d.Selected >= 0 && d.Selected < len(d.Options) {
		return d.Options[d.Selected].Value
	}
	return nil
}

// createToggleEffect creates effect when toggling dropdown
func (d *Dropdown) createToggleEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		d.Particles.AddSparkles(d.X+d.Width/2, d.Y, 6)
		return ParticleTickMsg{}
	})
}

// createSelectionEffect creates effect when selecting option
func (d *Dropdown) createSelectionEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(time.Time) tea.Msg {
		d.Particles.AddHearts(d.X+d.Width/2, d.Y+d.Selected+1, 4)
		return ParticleTickMsg{}
	})
}

// Render renders the stunning dropdown
func (d *Dropdown) Render() string {
	// Render header
	selectedText := "Select option..."
	if d.Selected >= 0 && d.Selected < len(d.Options) {
		selectedText = d.Options[d.Selected].Text
	}

	arrow := "▼"
	if d.Open {
		arrow = "▲"
	}

	headerText := fmt.Sprintf("%s %s", selectedText, arrow)

	style := d.Style
	if d.Focused {
		style = style.Copy().
			BorderForeground(charmtone.Coral).
			Background(lipgloss.Color("#fff8f8"))
	}

	header := style.Render(headerText)

	if !d.Open {
		return fmt.Sprintf("%s\n%s", d.Label, header)
	}

	// Render options
	var options []string
	for i, option := range d.Options {
		optStyle := d.OptionStyle
		if i == d.Selected {
			optStyle = d.SelectedStyle
		}
		options = append(options, optStyle.Render(option.Text))
	}

	optionsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Guac).
		Background(lipgloss.Color("#ffffff")).
		Render(strings.Join(options, "\n"))

	return fmt.Sprintf("%s\n%s\n%s", d.Label, header, optionsBox)
}

// Modal represents a modal dialog with overlay
type Modal struct {
	X, Y          int
	Width, Height int
	Title         string
	Content       string
	Buttons       []*Button
	Visible       bool
	Focused       bool
	Style         lipgloss.Style
	OverlayStyle  lipgloss.Style
	Particles     *ParticleSystem
	Animation     *AnimatedElement
}

// NewModal creates a stunning modal dialog
func NewModal(title, content string, width, height int) *Modal {
	modalStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.ThickBorder()).
		BorderForeground(charmtone.Salmon).
		Background(lipgloss.Color("#fff8f8")).
		Padding(2).
		Align(lipgloss.Center)

	overlayStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#00000080"))

	return &Modal{
		Width:        width,
		Height:       height,
		Title:        title,
		Content:      content,
		Buttons:      make([]*Button, 0),
		Style:        modalStyle,
		OverlayStyle: overlayStyle,
		Particles:    NewParticleSystem(width+20, height+20),
		Animation:    NewAnimatedElement(title, 0, 0),
	}
}

// AddButton adds a button to the modal
func (m *Modal) AddButton(button *Button) {
	m.Buttons = append(m.Buttons, button)
}

// Show shows the modal with animation
func (m *Modal) Show() {
	m.Visible = true
	m.Animation.Bounce(5, 0.8)
	m.Particles.AddSparkles(m.Width/2, m.Height/2, 15)
	m.Particles.AddHearts(m.Width/2, m.Height/2, 8)
}

// Hide hides the modal
func (m *Modal) Hide() {
	m.Visible = false
}

// Update updates the modal
func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	if !m.Visible {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		m.Particles.Update(0.05)
		m.Animation.Update(0.05)
		cmds = append(cmds, ParticleUpdateCmd())

	case tea.KeyMsg:
		if m.Focused {
			switch msg.String() {
			case "esc":
				m.Hide()
			case "tab", "shift+tab":
				// Handle button focus cycling
				// Implementation depends on specific needs
			}
		}

		// Update buttons
		for _, button := range m.Buttons {
			_, cmd := button.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// Focus sets modal focus
func (m *Modal) Focus() {
	m.Focused = true
}

// Blur removes modal focus
func (m *Modal) Blur() {
	m.Focused = false
}

// Render renders the stunning modal
func (m *Modal) Render() string {
	if !m.Visible {
		return ""
	}

	// Create title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(charmtone.Coral).
		Align(lipgloss.Center).
		Width(m.Width - 4)

	title := titleStyle.Render(m.Title)

	// Create content
	contentStyle := lipgloss.NewStyle().
		Width(m.Width - 4).
		Height(m.Height - 8).
		Align(lipgloss.Center)

	content := contentStyle.Render(m.Content)

	// Render buttons
	var buttonRow string
	if len(m.Buttons) > 0 {
		var buttons []string
		for _, button := range m.Buttons {
			buttons = append(buttons, button.Render())
		}
		buttonRow = lipgloss.JoinHorizontal(lipgloss.Center, buttons...)
	}

	// Combine all parts
	modalContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		content,
		"",
		buttonRow,
	)

	return m.Style.Render(modalContent)
}
