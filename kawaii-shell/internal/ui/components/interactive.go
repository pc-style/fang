package components

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// ButtonState represents the current state of a button
type ButtonState int

const (
	ButtonIdle ButtonState = iota
	ButtonHover
	ButtonPressed
	ButtonDisabled
)

// Button represents an interactive button with stunning effects
type Button struct {
	Text         string
	X, Y         int
	Width        int
	Height       int
	State        ButtonState
	Style        lipgloss.Style
	HoverStyle   lipgloss.Style
	PressedStyle lipgloss.Style
	OnClick      func()
	Particles    *ParticleSystem
	Animation    *AnimatedElement
	Focused      bool
	GlowLevel    float64
	PulseTime    float64
}

// NewButton creates a stunning new button
func NewButton(text string, x, y, width int) *Button {
	baseStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Coral).
		Background(lipgloss.Color("#fff8f8")).
		Foreground(charmtone.Charcoal).
		Bold(true)

	hoverStyle := baseStyle.Copy().
		BorderForeground(charmtone.Salmon).
		Background(lipgloss.Color("#ffe8e8")).
		Foreground(charmtone.Coral).
		Transform(func(s string) string {
			return fmt.Sprintf("✨ %s ✨", s)
		})

	pressedStyle := baseStyle.Copy().
		BorderForeground(charmtone.Cherry).
		Background(charmtone.Salmon).
		Foreground(charmtone.Butter).
		Transform(func(s string) string {
			return fmt.Sprintf("🌟 %s 🌟", s)
		})

	return &Button{
		Text:         text,
		X:            x,
		Y:            y,
		Width:        width,
		Height:       3,
		State:        ButtonIdle,
		Style:        baseStyle,
		HoverStyle:   hoverStyle,
		PressedStyle: pressedStyle,
		Particles:    NewParticleSystem(50, 20),
		Animation:    NewAnimatedElement(text, float64(x), float64(y)),
		GlowLevel:    0.0,
		PulseTime:    0.0,
	}
}

// Update updates the button state and animations
func (b *Button) Update(msg tea.Msg) (*Button, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		b.Particles.Update(0.05)
		b.updateAnimations()
		cmds = append(cmds, ParticleUpdateCmd())

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if b.Focused && b.State != ButtonDisabled {
				b.Press()
				cmds = append(cmds, b.createPressEffect())
			}
		case "tab", "shift+tab":
			// Handle focus changes in parent component
		}

	case tea.MouseMsg:
		if b.isMouseOver(msg.X, msg.Y) {
			switch msg.Type {
			case tea.MouseMotion:
				if b.State == ButtonIdle {
					b.Hover()
					cmds = append(cmds, b.createHoverEffect())
				}
			case tea.MouseLeft:
				b.Press()
				cmds = append(cmds, b.createPressEffect())
			}
		} else {
			if b.State == ButtonHover {
				b.Idle()
			}
		}
	}

	return b, tea.Batch(cmds...)
}

// Hover activates hover state with beautiful effects
func (b *Button) Hover() {
	if b.State == ButtonDisabled {
		return
	}
	b.State = ButtonHover
	b.Animation.Glow(0.8, 1.0)
	b.Particles.AddSparkles(b.X+b.Width/2, b.Y+b.Height/2, 5)
}

// Press activates pressed state with celebration
func (b *Button) Press() {
	if b.State == ButtonDisabled {
		return
	}
	b.State = ButtonPressed
	b.Animation.Bounce(3, 0.3)
	b.Particles.AddSparkles(b.X+b.Width/2, b.Y+b.Height/2, 15)
	b.Particles.AddHearts(b.X+b.Width/2, b.Y+b.Height/2, 8)

	if b.OnClick != nil {
		b.OnClick()
	}

	// Return to idle after press
	go func() {
		time.Sleep(200 * time.Millisecond)
		b.Idle()
	}()
}

// Idle returns button to idle state
func (b *Button) Idle() {
	if b.State == ButtonDisabled {
		return
	}
	b.State = ButtonIdle
	b.GlowLevel = 0.0
}

// Focus sets button focus
func (b *Button) Focus() {
	b.Focused = true
	b.Animation.Pulse(2.0)
}

// Blur removes button focus
func (b *Button) Blur() {
	b.Focused = false
}

// updateAnimations updates button animations
func (b *Button) updateAnimations() {
	b.PulseTime += 0.1

	// Update glow for focused buttons
	if b.Focused {
		b.GlowLevel = 0.3 + math.Sin(b.PulseTime)*0.2
	}

	// Update button animation
	b.Animation.Update(0.05)
}

// isMouseOver checks if mouse is over button
func (b *Button) isMouseOver(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

// createHoverEffect creates particle effect for hover
func (b *Button) createHoverEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		b.Particles.AddSparkles(b.X+b.Width/2, b.Y+b.Height/2, 2)
		return ParticleTickMsg{}
	})
}

// createPressEffect creates celebration effect for press
func (b *Button) createPressEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(time.Time) tea.Msg {
		b.Particles.AddSparkles(b.X+b.Width/2, b.Y+b.Height/2, 10)
		return ParticleTickMsg{}
	})
}

// Render renders the stunning button
func (b *Button) Render() string {
	var style lipgloss.Style

	switch b.State {
	case ButtonHover:
		style = b.HoverStyle
	case ButtonPressed:
		style = b.PressedStyle
	case ButtonDisabled:
		style = b.Style.Copy().
			Foreground(lipgloss.Color("#999999")).
			BorderForeground(lipgloss.Color("#cccccc"))
	default:
		style = b.Style
	}

	// Add focus indicator
	if b.Focused {
		style = style.Copy().
			Border(lipgloss.ThickBorder()).
			BorderForeground(charmtone.Malibu)
	}

	// Add glow effect
	if b.GlowLevel > 0 {
		glowIntensity := int(b.GlowLevel * 10)
		if glowIntensity > 0 {
			style = style.Copy().Border(lipgloss.ThickBorder())
		}
	}

	return style.Render(b.Text)
}

// Slider represents an interactive slider with particle trails
type Slider struct {
	X, Y      int
	Width     int
	Min, Max  float64
	Value     float64
	Step      float64
	Label     string
	Style     lipgloss.Style
	Particles *ParticleSystem
	Focused   bool
	Dragging  bool
	GlowPos   int
}

// NewSlider creates a stunning new slider
func NewSlider(label string, x, y, width int, min, max, value float64) *Slider {
	style := lipgloss.NewStyle().
		Width(width).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Guppy).
		Background(lipgloss.Color("#f0f8ff"))

	return &Slider{
		X:         x,
		Y:         y,
		Width:     width,
		Min:       min,
		Max:       max,
		Value:     value,
		Step:      1.0,
		Label:     label,
		Style:     style,
		Particles: NewParticleSystem(width+10, 10),
		GlowPos:   0,
	}
}

// Update updates the slider
func (s *Slider) Update(msg tea.Msg) (*Slider, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		s.Particles.Update(0.05)
		s.updateGlow()
		cmds = append(cmds, ParticleUpdateCmd())

	case tea.KeyMsg:
		if s.Focused {
			switch msg.String() {
			case "left", "h":
				s.setValue(s.Value - s.Step)
				s.createTrailEffect()
			case "right", "l":
				s.setValue(s.Value + s.Step)
				s.createTrailEffect()
			}
		}

	case tea.MouseMsg:
		if s.isMouseOver(msg.X, msg.Y) {
			switch msg.Type {
			case tea.MouseLeft:
				s.Dragging = true
				s.setValueFromMouse(msg.X)
				s.createTrailEffect()
			case tea.MouseMotion:
				if s.Dragging {
					s.setValueFromMouse(msg.X)
					s.createTrailEffect()
				}
			case tea.MouseRelease:
				s.Dragging = false
			}
		}
	}

	return s, tea.Batch(cmds...)
}

// setValue sets slider value with bounds checking
func (s *Slider) setValue(value float64) {
	if value < s.Min {
		value = s.Min
	}
	if value > s.Max {
		value = s.Max
	}
	s.Value = value
}

// setValueFromMouse sets value based on mouse position
func (s *Slider) setValueFromMouse(mouseX int) {
	relativeX := mouseX - s.X
	if relativeX < 0 {
		relativeX = 0
	}
	if relativeX > s.Width-4 {
		relativeX = s.Width - 4
	}

	percent := float64(relativeX) / float64(s.Width-4)
	s.setValue(s.Min + percent*(s.Max-s.Min))
}

// isMouseOver checks if mouse is over slider
func (s *Slider) isMouseOver(x, y int) bool {
	return x >= s.X && x <= s.X+s.Width && y >= s.Y && y <= s.Y+3
}

// updateGlow updates the glow position
func (s *Slider) updateGlow() {
	percent := (s.Value - s.Min) / (s.Max - s.Min)
	s.GlowPos = int(percent * float64(s.Width-4))
}

// createTrailEffect creates particle trail when slider moves
func (s *Slider) createTrailEffect() {
	x := s.X + s.GlowPos + 2
	y := s.Y + 1
	s.Particles.AddSparkles(x, y, 3)
}

// Focus sets slider focus
func (s *Slider) Focus() {
	s.Focused = true
}

// Blur removes slider focus
func (s *Slider) Blur() {
	s.Focused = false
}

// Render renders the stunning slider
func (s *Slider) Render() string {
	percent := (s.Value - s.Min) / (s.Max - s.Min)
	filledWidth := int(percent * float64(s.Width-4))

	// Create the slider track
	_ = strings.Repeat("─", s.Width-4)
	filled := strings.Repeat("━", filledWidth)
	empty := strings.Repeat("─", s.Width-4-filledWidth)

	// Create the slider handle
	handle := "●"
	if s.Focused {
		handle = "◉"
	}
	if s.Dragging {
		handle = "✨"
	}

	// Position the handle
	var sliderBar string
	if filledWidth > 0 {
		sliderBar = filled[:filledWidth] + handle
		if filledWidth < s.Width-4 {
			sliderBar += empty
		}
	} else {
		sliderBar = handle + empty
	}

	// Create the complete slider
	content := fmt.Sprintf("%s\n├%s┤\n%.1f", s.Label, sliderBar, s.Value)

	style := s.Style
	if s.Focused {
		style = style.Copy().
			BorderForeground(charmtone.Coral).
			Background(lipgloss.Color("#fff8f8"))
	}

	return style.Render(content)
}

// ProgressBar represents an animated progress bar with rainbow effects
type ProgressBar struct {
	X, Y       int
	Width      int
	Progress   float64
	Max        float64
	Label      string
	Style      lipgloss.Style
	Particles  *ParticleSystem
	RainbowPos int
	PulseTime  float64
}

// NewProgressBar creates a stunning progress bar
func NewProgressBar(label string, x, y, width int, max float64) *ProgressBar {
	style := lipgloss.NewStyle().
		Width(width).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Pony).
		Background(lipgloss.Color("#fff0ff"))

	return &ProgressBar{
		X:         x,
		Y:         y,
		Width:     width,
		Max:       max,
		Label:     label,
		Style:     style,
		Particles: NewParticleSystem(width+10, 10),
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(msg tea.Msg) (*ProgressBar, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ParticleTickMsg:
		pb.Particles.Update(0.05)
		pb.updateRainbow()
		pb.PulseTime += 0.1

		// Add sparkles at the progress edge
		if pb.Progress > 0 && pb.Progress < pb.Max {
			percent := pb.Progress / pb.Max
			x := pb.X + int(percent*float64(pb.Width-4)) + 2
			y := pb.Y + 1
			pb.Particles.AddSparkles(x, y, 1)
		}

		cmds = append(cmds, ParticleUpdateCmd())
	}

	return pb, tea.Batch(cmds...)
}

// SetProgress sets the progress value
func (pb *ProgressBar) SetProgress(progress float64) {
	if progress > pb.Max {
		progress = pb.Max
	}
	if progress < 0 {
		progress = 0
	}

	oldProgress := pb.Progress
	pb.Progress = progress

	// Create celebration effect when progress increases
	if progress > oldProgress {
		percent := progress / pb.Max
		x := pb.X + int(percent*float64(pb.Width-4)) + 2
		y := pb.Y + 1
		pb.Particles.AddSparkles(x, y, 5)

		// Extra celebration when complete
		if progress >= pb.Max {
			pb.Particles.AddHearts(pb.X+pb.Width/2, pb.Y+1, 10)
			pb.Particles.AddFlowerPetals(pb.X+pb.Width/2, pb.Y+1, 8)
		}
	}
}

// updateRainbow updates rainbow animation
func (pb *ProgressBar) updateRainbow() {
	pb.RainbowPos = (pb.RainbowPos + 1) % 7
}

// Render renders the stunning progress bar
func (pb *ProgressBar) Render() string {
	percent := pb.Progress / pb.Max
	if percent > 1.0 {
		percent = 1.0
	}

	filledWidth := int(percent * float64(pb.Width-4))
	emptyWidth := (pb.Width - 4) - filledWidth

	// Create rainbow effect for filled portion
	colors := []string{
		"#ff0000", // Red
		"#ff8000", // Orange
		"#ffff00", // Yellow
		"#00ff00", // Green
		"#0080ff", // Blue
		"#8000ff", // Purple
		"#ff00ff", // Magenta
	}

	var filled string
	for i := 0; i < filledWidth; i++ {
		colorIndex := (i + pb.RainbowPos) % len(colors)
		char := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors[colorIndex])).
			Bold(true).
			Render("█")
		filled += char
	}

	empty := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cccccc")).
		Render(strings.Repeat("░", emptyWidth))

	// Add pulse effect when complete
	percentText := fmt.Sprintf("%.0f%%", percent*100)
	if percent >= 1.0 {
		pulseIntensity := math.Sin(pb.PulseTime)*0.5 + 0.5
		if pulseIntensity > 0.7 {
			percentText = fmt.Sprintf("✨ %s ✨", percentText)
		}
	}

	content := fmt.Sprintf("%s\n[%s%s] %s", pb.Label, filled, empty, percentText)

	return pb.Style.Render(content)
}
