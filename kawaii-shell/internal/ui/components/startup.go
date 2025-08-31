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

// StartupPhase represents different phases of the startup sequence
type StartupPhase int

const (
	PhaseInitial StartupPhase = iota
	PhaseLogoReveal
	PhaseTextMorphing
	PhaseParticleExplosion
	PhaseInfoCascade
	PhaseFinalReveal
	PhaseComplete
)

// StartupSequence manages the stunning startup animation
type StartupSequence struct {
	phase             StartupPhase
	startTime         time.Time
	phaseStartTime    time.Time
	animationManager  *AnimationManager
	particleSystem    *ParticleSystem
	currentFrame      int
	totalFrames       int
	logoAlpha         float64
	textMorphProgress float64
	cascadeProgress   float64
	explosionRadius   float64
	glowIntensity     float64
	sparkleCount      int
	rainbowOffset     int
	pulseAnimation    float64
	width             int
	height            int
	completed         bool
	currentLogo       []string
	finalLogo         []string
	infoLines         []string
	cascadeDelay      []float64
	version           string
}

// NewStartupSequence creates a stunning startup animation
func NewStartupSequence(width, height int, version string) *StartupSequence {
	finalLogo := []string{
		"   ╭─────────────────────────────────────╮",
		"   │  🌸✨ KAWAII SHELL ✨🌸            │",
		"   │                                     │",
		"   │     ／| ／| 　♡   Stunning Terminal  │",
		"   │    (  ˘ ᵕ ˘ )   Experience  ⭐     │",
		"   │     ○_○━━━━━━━━━━━━━━━━━━━○_○         │",
		"   │                                     │",
		"   │  🎀 Making terminals magical! 🎀   │",
		"   ╰─────────────────────────────────────╯",
	}

	infoLines := []string{
		"🚀 Initializing stunning visual effects...",
		"✨ Loading particle systems...",
		"🎨 Applying gorgeous themes...",
		"🤖 Awakening AI pet companion...",
		"🌈 Calibrating rainbow generators...",
		"💖 Spreading kawaii energy...",
		"🎪 Ready for magical adventures!",
	}

	// Create cascade delays for staggered animation
	cascadeDelay := make([]float64, len(infoLines))
	for i := range cascadeDelay {
		cascadeDelay[i] = float64(i) * 0.3
	}

	return &StartupSequence{
		phase:            PhaseInitial,
		startTime:        time.Now(),
		phaseStartTime:   time.Now(),
		animationManager: NewAnimationManager(),
		particleSystem:   NewParticleSystem(width, height),
		totalFrames:      60,
		width:            width,
		height:           height,
		finalLogo:        finalLogo,
		infoLines:        infoLines,
		cascadeDelay:     cascadeDelay,
		version:          version,
		currentLogo:      make([]string, len(finalLogo)),
	}
}

// StartupTickMsg represents animation tick for startup
type StartupTickMsg struct {
	Time time.Time
}

// StartupCompleteMsg signals startup completion
type StartupCompleteMsg struct{}

// Init initializes the startup sequence
func (ss *StartupSequence) Init() tea.Cmd {
	return tea.Batch(
		ParticleUpdateCmd(),
		tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return StartupTickMsg{Time: t}
		}),
	)
}

// Update updates the startup sequence with cinematic timing
func (ss *StartupSequence) Update(msg tea.Msg) (*StartupSequence, tea.Cmd) {
	if ss.completed {
		return ss, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case StartupTickMsg:
		ss.updateAnimations()
		ss.updatePhase()
		cmds = append(cmds, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return StartupTickMsg{Time: t}
		}))

	case ParticleTickMsg:
		ss.particleSystem.Update(0.05)
		ss.animationManager.Update()
		cmds = append(cmds, ParticleUpdateCmd())
	}

	return ss, tea.Batch(cmds...)
}

// updateAnimations updates all animation values
func (ss *StartupSequence) updateAnimations() {
	ss.currentFrame++
	elapsed := time.Since(ss.phaseStartTime).Seconds()
	totalElapsed := time.Since(ss.startTime).Seconds()

	// Update global animation values
	ss.pulseAnimation = math.Sin(totalElapsed*3)*0.3 + 0.7
	ss.rainbowOffset = (ss.rainbowOffset + 1) % 7

	switch ss.phase {
	case PhaseInitial:
		// Fade in effect
		ss.logoAlpha = math.Min(elapsed/1.0, 1.0)

	case PhaseLogoReveal:
		// Cascading logo reveal with bounce
		progress := math.Min(elapsed/2.0, 1.0)
		bounceProgress := EaseBounce(progress)
		ss.logoAlpha = bounceProgress

		// Create sparkle effects during reveal
		if ss.currentFrame%5 == 0 {
			ss.particleSystem.AddSparkles(ss.width/2, ss.height/3, 3)
		}

	case PhaseTextMorphing:
		// Smooth text morphing animation
		ss.textMorphProgress = EaseInOutQuad(math.Min(elapsed/1.5, 1.0))

		// Morphing particle effects
		if ss.currentFrame%3 == 0 {
			ss.particleSystem.AddSparkles(ss.width/2, ss.height/2, 2)
		}

	case PhaseParticleExplosion:
		// Dramatic particle explosion
		progress := math.Min(elapsed/2.0, 1.0)
		ss.explosionRadius = EaseOutQuad(progress) * float64(ss.width/3)
		ss.glowIntensity = math.Sin(progress*math.Pi*2)*0.5 + 0.5

		// Create explosion particles
		if progress < 0.8 && ss.currentFrame%2 == 0 {
			angle := float64(ss.currentFrame) * 0.5
			x := ss.width/2 + int(math.Cos(angle)*ss.explosionRadius*0.5)
			y := ss.height/2 + int(math.Sin(angle)*ss.explosionRadius*0.3)
			ss.particleSystem.AddSparkles(x, y, 5)
			ss.particleSystem.AddHearts(x, y, 2)
		}

	case PhaseInfoCascade:
		// Cascading information reveal
		ss.cascadeProgress = EaseOutQuad(math.Min(elapsed/3.0, 1.0))

		// Add cascade particles
		if ss.currentFrame%4 == 0 {
			for i, delay := range ss.cascadeDelay {
				if elapsed > delay {
					y := ss.height/2 + i*2 + 10
					ss.particleSystem.AddSparkles(ss.width/4, y, 1)
				}
			}
		}

	case PhaseFinalReveal:
		// Final reveal with maximum effect
		progress := math.Min(elapsed/2.0, 1.0)
		ss.logoAlpha = 1.0
		ss.glowIntensity = progress
		ss.sparkleCount = int(progress * 20)

		// Continuous sparkle shower
		if ss.currentFrame%2 == 0 {
			for i := 0; i < 3; i++ {
				x := ss.width/4 + i*(ss.width/2)
				ss.particleSystem.AddSparkles(x, ss.height/6, 2)
			}
		}
	}
}

// updatePhase manages phase transitions
func (ss *StartupSequence) updatePhase() {
	elapsed := time.Since(ss.phaseStartTime).Seconds()

	var nextPhase StartupPhase
	var phaseDuration float64

	switch ss.phase {
	case PhaseInitial:
		nextPhase = PhaseLogoReveal
		phaseDuration = 1.0
	case PhaseLogoReveal:
		nextPhase = PhaseTextMorphing
		phaseDuration = 2.5
	case PhaseTextMorphing:
		nextPhase = PhaseParticleExplosion
		phaseDuration = 2.0
	case PhaseParticleExplosion:
		nextPhase = PhaseInfoCascade
		phaseDuration = 2.5
	case PhaseInfoCascade:
		nextPhase = PhaseFinalReveal
		phaseDuration = 3.5
	case PhaseFinalReveal:
		nextPhase = PhaseComplete
		phaseDuration = 2.0
	case PhaseComplete:
		ss.completed = true
		return
	}

	if elapsed >= phaseDuration {
		ss.phase = nextPhase
		ss.phaseStartTime = time.Now()
		ss.createPhaseTransitionEffect()
	}
}

// createPhaseTransitionEffect creates special effects during phase transitions
func (ss *StartupSequence) createPhaseTransitionEffect() {
	centerX := ss.width / 2
	centerY := ss.height / 2

	switch ss.phase {
	case PhaseLogoReveal:
		// Explosion of sparkles when logo appears
		ss.particleSystem.AddSparkles(centerX, centerY, 15)

	case PhaseTextMorphing:
		// Hearts explosion for text morphing
		ss.particleSystem.AddHearts(centerX, centerY, 10)

	case PhaseParticleExplosion:
		// Massive particle burst
		ss.particleSystem.AddSparkles(centerX, centerY, 25)
		ss.particleSystem.AddFlowerPetals(centerX, centerY, 12)

	case PhaseInfoCascade:
		// Wave of sparkles across screen
		for i := 0; i < ss.width/10; i++ {
			x := i * 10
			ss.particleSystem.AddSparkles(x, centerY, 2)
		}

	case PhaseFinalReveal:
		// Celebration explosion
		ss.particleSystem.AddSparkles(centerX, centerY, 30)
		ss.particleSystem.AddHearts(centerX, centerY, 15)
		ss.particleSystem.AddFlowerPetals(centerX, centerY, 20)
	}
}

// IsComplete returns whether startup is finished
func (ss *StartupSequence) IsComplete() bool {
	return ss.completed
}

// GetDuration returns total startup duration
func (ss *StartupSequence) GetDuration() time.Duration {
	return time.Second * 14 // Total duration across all phases
}

// Render renders the stunning startup sequence
func (ss *StartupSequence) Render() string {
	if ss.completed {
		return ""
	}

	var content strings.Builder

	// Create the main container with stunning border
	containerStyle := lipgloss.NewStyle().
		Width(ss.width).
		Height(ss.height).
		Border(lipgloss.ThickBorder()).
		BorderForeground(charmtone.Coral).
		Background(lipgloss.Color("#0a0a1a")).
		Padding(2).
		Align(lipgloss.Center)

	// Add glow effect based on phase
	if ss.glowIntensity > 0 {
		containerStyle = containerStyle.
			BorderForeground(ss.getGlowColor())
	}

	// Render based on current phase
	switch ss.phase {
	case PhaseInitial, PhaseLogoReveal:
		content.WriteString(ss.renderLogoReveal())

	case PhaseTextMorphing:
		content.WriteString(ss.renderTextMorphing())

	case PhaseParticleExplosion:
		content.WriteString(ss.renderParticleExplosion())

	case PhaseInfoCascade:
		content.WriteString(ss.renderInfoCascade())

	case PhaseFinalReveal:
		content.WriteString(ss.renderFinalReveal())
	}

	// Add particle overlay
	particles := ss.renderParticleOverlay()
	if particles != "" {
		content.WriteString("\n\n")
		content.WriteString(particles)
	}

	return containerStyle.Render(content.String())
}

// renderLogoReveal renders the logo reveal phase
func (ss *StartupSequence) renderLogoReveal() string {
	var result strings.Builder

	// Title with stunning effects
	titleStyle := lipgloss.NewStyle().
		Foreground(charmtone.Coral).
		Bold(true).
		Align(lipgloss.Center).
		Transform(func(s string) string {
			return ss.applyRainbowEffect(s)
		})

	if ss.logoAlpha > 0 {
		alpha := int(ss.logoAlpha * 100)
		title := fmt.Sprintf("🌸✨ KAWAII SHELL v%s ✨🌸", ss.version)

		if alpha < 100 {
			// Fade in effect
			fadeStyle := titleStyle.Copy().
				Foreground(lipgloss.Color(fmt.Sprintf("#ff%02x%02x", alpha*2, alpha*2)))
			result.WriteString(fadeStyle.Render(title))
		} else {
			result.WriteString(titleStyle.Render(title))
		}
	}

	// Subtitle
	if ss.logoAlpha > 0.5 {
		result.WriteString("\n\n")
		subtitleStyle := lipgloss.NewStyle().
			Foreground(charmtone.Guppy).
			Italic(true).
			Align(lipgloss.Center)

		result.WriteString(subtitleStyle.Render("Making terminals magical! ✨"))
	}

	return result.String()
}

// renderTextMorphing renders the text morphing phase
func (ss *StartupSequence) renderTextMorphing() string {
	var result strings.Builder

	// Morphing logo effect
	for i, line := range ss.finalLogo {
		if ss.textMorphProgress > float64(i)/float64(len(ss.finalLogo)) {
			// Apply morphing effects
			morphedLine := ss.applyMorphingEffect(line, ss.textMorphProgress)

			lineStyle := lipgloss.NewStyle().
				Foreground(ss.getMorphColor(i)).
				Bold(i == 1 || i == 7). // Highlight title and footer
				Align(lipgloss.Center)

			result.WriteString(lineStyle.Render(morphedLine))
			result.WriteString("\n")
		}
	}

	return result.String()
}

// renderParticleExplosion renders the particle explosion phase
func (ss *StartupSequence) renderParticleExplosion() string {
	var result strings.Builder

	// Render logo with explosion effects
	for i, line := range ss.finalLogo {
		explosionIntensity := math.Max(0, ss.explosionRadius/float64(ss.width)*2-float64(i)*0.1)

		lineStyle := lipgloss.NewStyle().
			Foreground(ss.getExplosionColor(explosionIntensity)).
			Bold(true).
			Align(lipgloss.Center)

		if explosionIntensity > 0.5 {
			// Add transformation for high intensity
			lineStyle = lineStyle.Transform(func(s string) string {
				return fmt.Sprintf("✨ %s ✨", s)
			})
		}

		result.WriteString(lineStyle.Render(line))
		result.WriteString("\n")
	}

	return result.String()
}

// renderInfoCascade renders the information cascade phase
func (ss *StartupSequence) renderInfoCascade() string {
	var result strings.Builder

	// Render logo
	logoStyle := lipgloss.NewStyle().
		Foreground(charmtone.Coral).
		Bold(true).
		Align(lipgloss.Center)

	for _, line := range ss.finalLogo {
		result.WriteString(logoStyle.Render(line))
		result.WriteString("\n")
	}

	result.WriteString("\n\n")

	// Render cascading info lines
	elapsed := time.Since(ss.phaseStartTime).Seconds()

	for i, info := range ss.infoLines {
		if elapsed > ss.cascadeDelay[i] {
			// Calculate fade-in progress
			progress := math.Min((elapsed-ss.cascadeDelay[i])/0.5, 1.0)

			infoStyle := lipgloss.NewStyle().
				Foreground(ss.getCascadeColor(i, progress)).
				Bold(progress > 0.8).
				Margin(0, 2)

			if progress > 0.9 {
				infoStyle = infoStyle.Transform(func(s string) string {
					return fmt.Sprintf("  %s", s)
				})
			}

			result.WriteString(infoStyle.Render(info))
			result.WriteString("\n")
		}
	}

	return result.String()
}

// renderFinalReveal renders the final reveal phase
func (ss *StartupSequence) renderFinalReveal() string {
	var result strings.Builder

	// Full logo with maximum effects
	for i, line := range ss.finalLogo {
		lineStyle := lipgloss.NewStyle().
			Foreground(ss.getFinalColor(i)).
			Bold(true).
			Align(lipgloss.Center).
			Transform(func(s string) string {
				if i == 1 || i == 7 { // Title and footer
					return ss.applyRainbowEffect(s)
				}
				return s
			})

		result.WriteString(lineStyle.Render(line))
		result.WriteString("\n")
	}

	result.WriteString("\n")

	// Final message
	finalStyle := lipgloss.NewStyle().
		Foreground(charmtone.Malibu).
		Bold(true).
		Align(lipgloss.Center).
		Transform(func(s string) string {
			return fmt.Sprintf("🎉 %s 🎉", s)
		})

	result.WriteString(finalStyle.Render("Welcome to your magical terminal!"))

	return result.String()
}

// renderParticleOverlay renders particle effects overlay
func (ss *StartupSequence) renderParticleOverlay() string {
	// This would integrate with the particle system
	// For now, return empty as particles are handled by the system
	return ""
}

// Helper functions for stunning effects

func (ss *StartupSequence) applyRainbowEffect(text string) string {
	colors := []string{
		"#ff0000", // Red
		"#ff8000", // Orange
		"#ffff00", // Yellow
		"#00ff00", // Green
		"#0080ff", // Blue
		"#8000ff", // Purple
		"#ff00ff", // Magenta
	}

	var result strings.Builder
	for i, char := range text {
		if char == ' ' {
			result.WriteRune(char)
			continue
		}

		colorIndex := (i + ss.rainbowOffset) % len(colors)
		styledChar := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors[colorIndex])).
			Bold(true).
			Render(string(char))
		result.WriteString(styledChar)
	}

	return result.String()
}

func (ss *StartupSequence) applyMorphingEffect(text string, progress float64) string {
	// Simple morphing effect - could be enhanced with character substitution
	intensity := int(progress * 100)
	if intensity < 50 {
		// Partial reveal
		visible := int(float64(len(text)) * progress * 2)
		if visible > len(text) {
			visible = len(text)
		}
		return text[:visible]
	}
	return text
}

func (ss *StartupSequence) getGlowColor() lipgloss.Color {
	intensity := int(ss.glowIntensity * 255)
	return lipgloss.Color(fmt.Sprintf("#%02xff%02x", intensity, intensity))
}

func (ss *StartupSequence) getMorphColor(lineIndex int) lipgloss.Color {
	colors := []lipgloss.Color{
		charmtone.Coral,
		charmtone.Salmon,
		charmtone.Guppy,
		charmtone.Malibu,
		charmtone.Pony,
	}
	return colors[lineIndex%len(colors)]
}

func (ss *StartupSequence) getExplosionColor(intensity float64) lipgloss.Color {
	if intensity > 0.8 {
		return lipgloss.Color("#ffffff")
	} else if intensity > 0.5 {
		return lipgloss.Color("#ffff00")
	} else if intensity > 0.2 {
		return lipgloss.Color("#ff8000")
	}
	return charmtone.Coral
}

func (ss *StartupSequence) getCascadeColor(index int, progress float64) lipgloss.Color {
	baseColors := []lipgloss.Color{
		charmtone.Coral,
		charmtone.Salmon,
		charmtone.Guppy,
		charmtone.Malibu,
		charmtone.Pony,
		charmtone.Cherry,
		charmtone.Butter,
	}

	baseColor := baseColors[index%len(baseColors)]

	// Fade effect based on progress
	if progress < 1.0 {
		alpha := int(progress * 255)
		// Simple fade approximation
		return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", alpha, alpha, alpha))
	}

	return baseColor
}

func (ss *StartupSequence) getFinalColor(lineIndex int) lipgloss.Color {
	// Cycle through gorgeous colors with pulse effect
	colors := []lipgloss.Color{
		charmtone.Coral,
		charmtone.Salmon,
		charmtone.Guppy,
		charmtone.Malibu,
		charmtone.Pony,
	}

	baseIndex := (lineIndex + int(ss.pulseAnimation*10)) % len(colors)
	return colors[baseIndex]
}

// Easing functions for smooth animations
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}
