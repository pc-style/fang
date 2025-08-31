package themes

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// KawaiiTheme represents a ultra-cute theme with stunning visuals
type KawaiiTheme struct {
	Name           string
	Styles         KawaiiStyles
	GradientColors []string
	AnimationTime  time.Time
}

// KawaiiStyles contains all the stunning kawaii styling
type KawaiiStyles struct {
	// Basic text styles with gradient support
	Prompt    lipgloss.Style
	Input     lipgloss.Style
	Cursor    lipgloss.Style
	OutputBox lipgloss.Style
	InputBox  lipgloss.Style

	// Special content styles with effects
	CommandInfo lipgloss.Style
	Warning     lipgloss.Style
	Help        lipgloss.Style
	Info        lipgloss.Style
	Pet         lipgloss.Style
	PetBox      lipgloss.Style

	// New stunning styles
	Title       lipgloss.Style
	Success     lipgloss.Style
	Error       lipgloss.Style
	Sparkle     lipgloss.Style
	Highlight   lipgloss.Style
	Glow        lipgloss.Style
	Rainbow     lipgloss.Style
	FloatingBox lipgloss.Style
}

// Create animated gradient colors
func (kt *KawaiiTheme) GetAnimatedGradient() []string {
	if len(kt.GradientColors) < 2 {
		return kt.GradientColors
	}
	elapsed := time.Since(kt.AnimationTime).Seconds()
	_ = math.Sin(elapsed*0.5)*0.5 + 0.5
	return kt.GradientColors
}

// NewSakuraTheme creates the most beautiful sakura theme ever
func NewSakuraTheme() *KawaiiTheme {
	gradientColors := []string{
		string(charmtone.Coral),
		string(charmtone.Salmon),
		string(charmtone.Cherry),
		string(charmtone.Pony),
	}

	return &KawaiiTheme{
		Name:           "Sakura Dreams 🌸✨",
		GradientColors: gradientColors,
		AnimationTime:  time.Now(),
		Styles: KawaiiStyles{
			// Enhanced prompt with glow effect
			Prompt: lipgloss.NewStyle().
				Foreground(charmtone.Coral).
				Background(lipgloss.Color("#1a0a0a")).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(charmtone.Salmon).
				MarginRight(1),

			Input: lipgloss.NewStyle().
				Foreground(charmtone.Charcoal).
				Background(lipgloss.Color("#fff8f8")).
				Padding(0, 1).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(charmtone.Coral),

			Cursor: lipgloss.NewStyle().
				Foreground(charmtone.Butter).
				Background(charmtone.Coral).
				Bold(true).
				Blink(true),

			// Stunning output box with shadow effect
			OutputBox: lipgloss.NewStyle().
				Padding(2, 3).
				Border(lipgloss.ThickBorder()).
				BorderForeground(charmtone.Salmon).
				Background(lipgloss.Color("#fff5f5")).
				Foreground(charmtone.Charcoal).
				MarginBottom(1).

			InputBox: lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(charmtone.Coral).
				Background(lipgloss.Color("#ffe8e8")).

			CommandInfo: lipgloss.NewStyle().
				Foreground(charmtone.Malibu).
				Background(lipgloss.Color("#e8f4ff")).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(charmtone.Guppy).
				MarginBottom(1),

			Warning: lipgloss.NewStyle().
				Foreground(charmtone.Butter).
				Background(charmtone.Cherry).
				Padding(1, 2).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#ff6b6b")).
				Bold(true).
				Blink(true),

			Help: lipgloss.NewStyle().
				Foreground(charmtone.Cheeky).
				Background(lipgloss.Color("#f0fff0")).
				Padding(0, 1).
				Italic(true),

			Info: lipgloss.NewStyle().
				Foreground(charmtone.Pony).
				Background(lipgloss.Color("#fff0ff")).
				Padding(0, 1),

			Pet: lipgloss.NewStyle().
				Foreground(charmtone.Guac).
				Background(lipgloss.Color("#f5fff5")).
				Bold(true).
				Padding(0, 1),

			// Gorgeous pet box with floating effect
			PetBox: lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(charmtone.Salmon).
				Background(lipgloss.Color("#fff0f8")).
				Padding(1, 2).
				Align(lipgloss.Center).
				MarginTop(1),

			// New stunning styles
			Title: lipgloss.NewStyle().
				Foreground(charmtone.Coral).
				Background(lipgloss.Color("#2a1a1a")).
				Bold(true).
				Padding(1, 3).
				Border(lipgloss.ThickBorder()).
				BorderForeground(charmtone.Salmon).
				Align(lipgloss.Center).
				MarginBottom(2),

			Success: lipgloss.NewStyle().
				Foreground(charmtone.Butter).
				Background(charmtone.Cheeky).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(charmtone.Guac),

			Error: lipgloss.NewStyle().
				Foreground(charmtone.Butter).
				Background(charmtone.Cherry).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#cc0000")),

			Sparkle: lipgloss.NewStyle().
				Foreground(charmtone.Butter).
				Background(charmtone.Pony).
				Bold(true).
				Blink(true),

			Highlight: lipgloss.NewStyle().
				Foreground(charmtone.Charcoal).
				Background(charmtone.Butter).
				Bold(true).
				Padding(0, 1),

			Glow: lipgloss.NewStyle().
				Foreground(charmtone.Salmon).
				Bold(true).
				Italic(true),

			Rainbow: lipgloss.NewStyle().
				Bold(true).
				Foreground(charmtone.Malibu),

			FloatingBox: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(charmtone.Coral).
				Background(lipgloss.Color("#fff8f8")).
				Padding(1, 2).
				Align(lipgloss.Center),
		},
	}
}

// NewGalaxyTheme creates a stunning space/galaxy theme
func NewGalaxyTheme() *KawaiiTheme {
	gradientColors := []string{
		"#1a0033",
		"#330066",
		"#4d0080",
		"#6600cc",
		"#8000ff",
	}

	return &KawaiiTheme{
		Name:           "Galaxy Dreams 🌌⭐",
		GradientColors: gradientColors,
		AnimationTime:  time.Now(),
		Styles: KawaiiStyles{
			Prompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff66ff")).
				Background(lipgloss.Color("#1a0033")).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#8000ff")).
				MarginRight(1),

			Input: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e6ccff")).
				Background(lipgloss.Color("#0d001a")).
				Padding(0, 1),

			Cursor: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffff00")).
				Background(lipgloss.Color("#8000ff")).
				Bold(true).
				Blink(true),

			OutputBox: lipgloss.NewStyle().
				Padding(2, 3).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#6600cc")).
				Background(lipgloss.Color("#0d001a")).
				Foreground(lipgloss.Color("#ccccff")).
				MarginBottom(1).

			InputBox: lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#ff66ff")).
				Background(lipgloss.Color("#1a0033")).

			CommandInfo: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Background(lipgloss.Color("#001a1a")).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#0099cc")).
				MarginBottom(1),

			Warning: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#ff3333")).
				Padding(1, 2).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#cc0000")).
				Bold(true).
				Blink(true),

			Help: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#99ffcc")).
				Background(lipgloss.Color("#001a0d")).
				Padding(0, 1).
				Italic(true),

			Pet: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffcc99")).
				Background(lipgloss.Color("#1a0d00")).
				Bold(true).
				Padding(0, 1),

			PetBox: lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#ff66ff")).
				Background(lipgloss.Color("#1a0033")).
				Padding(1, 2).
				Align(lipgloss.Center).
				MarginTop(1),

			Title: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff66ff")).
				Background(lipgloss.Color("#1a0033")).
				Bold(true).
				Padding(1, 3).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#8000ff")).
				Align(lipgloss.Center).
				MarginBottom(2),

			Success: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#00ff66")).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#00cc33")),

			Sparkle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffff00")).
				Background(lipgloss.Color("#8000ff")).
				Bold(true).
				Blink(true),

			Glow: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff66ff")).
				Bold(true).
				Italic(true),
		},
	}
}

// NewCyberTheme creates a cyberpunk theme
func NewCyberTheme() *KawaiiTheme {
	gradientColors := []string{
		"#ff0080",
		"#00ff80",
		"#0080ff",
		"#8000ff",
	}

	return &KawaiiTheme{
		Name:           "Cyber Kawaii 🤖💫",
		GradientColors: gradientColors,
		AnimationTime:  time.Now(),
		Styles: KawaiiStyles{
			Prompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff80")).
				Background(lipgloss.Color("#001a00")).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#00cc66")).
				MarginRight(1),

			Input: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ccffcc")).
				Background(lipgloss.Color("#000000")).
				Padding(0, 1),

			OutputBox: lipgloss.NewStyle().
				Padding(2, 3).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#ff0080")).
				Background(lipgloss.Color("#000000")).
				Foreground(lipgloss.Color("#00ff80")).
				MarginBottom(1),

			CommandInfo: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0080ff")).
				Background(lipgloss.Color("#000019")).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#0066cc")),

			Pet: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0080")).
				Bold(true),

			PetBox: lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#00ff80")).
				Background(lipgloss.Color("#001a00")).
				Padding(1, 2).
				Align(lipgloss.Center),

			Sparkle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#ff0080")).
				Bold(true).
				Blink(true),
		},
	}
}

// NewOceanTheme creates an enhanced ocean theme
func NewOceanTheme() *KawaiiTheme {
	gradientColors := []string{
		string(charmtone.Malibu),
		string(charmtone.Guppy),
		"#0066cc",
		"#004499",
	}

	return &KawaiiTheme{
		Name:           "Ocean Breeze 🌊🐚",
		GradientColors: gradientColors,
		AnimationTime:  time.Now(),
		Styles: KawaiiStyles{
			Prompt: lipgloss.NewStyle().
				Foreground(charmtone.Malibu).
				Background(lipgloss.Color("#001a33")).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(charmtone.Guppy).
				MarginRight(1),

			Input: lipgloss.NewStyle().
				Foreground(charmtone.Charcoal).
				Background(lipgloss.Color("#f0f8ff")).
				Padding(0, 1),

			OutputBox: lipgloss.NewStyle().
				Padding(2, 3).
				Border(lipgloss.ThickBorder()).
				BorderForeground(charmtone.Guppy).
				Background(lipgloss.Color("#f0f8ff")).
				Foreground(charmtone.Charcoal).

			Pet: lipgloss.NewStyle().
				Foreground(charmtone.Guac).
				Bold(true),

			PetBox: lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(charmtone.Guppy).
				Background(lipgloss.Color("#e6f3ff")).
				Padding(1, 2).
				Align(lipgloss.Center).
		},
	}
}

// NewRainbowTheme creates the ultimate rainbow theme
func NewRainbowTheme() *KawaiiTheme {
	gradientColors := []string{
		"#ff0000", // Red
		"#ff8000", // Orange
		"#ffff00", // Yellow
		"#00ff00", // Green
		"#0000ff", // Blue
		"#8000ff", // Purple
		"#ff00ff", // Magenta
	}

	return &KawaiiTheme{
		Name:           "Rainbow Magic 🌈✨",
		GradientColors: gradientColors,
		AnimationTime:  time.Now(),
		Styles: KawaiiStyles{
			Prompt: lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				MarginRight(1),

			OutputBox: lipgloss.NewStyle().
				Padding(2, 3).
				Border(lipgloss.ThickBorder()).
				Background(lipgloss.Color("#fefefe")).

			Rainbow: lipgloss.NewStyle().
				Bold(true),

			Sparkle: lipgloss.NewStyle().
				Bold(true).
				Blink(true),
		},
	}
}

// ApplyRainbowEffect applies rainbow colors to text
func (kt *KawaiiTheme) ApplyRainbowEffect(text string) string {
	if kt.Name != "Rainbow Magic 🌈✨" {
		return text
	}

	colors := []string{
		"#ff0000", // Red
		"#ff8000", // Orange
		"#ffff00", // Yellow
		"#00ff00", // Green
		"#0080ff", // Blue
		"#8000ff", // Purple
	}

	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		styledChar := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).Render(string(char))
		result += styledChar
	}
	return result
}

// CreateSparkleText creates sparkling animated text
func (kt *KawaiiTheme) CreateSparkleText(text string) string {
	sparkles := []string{"✨", "⭐", "💫", "🌟", "⚡", "💎"}
	sparkle1 := sparkles[int(time.Since(kt.AnimationTime).Seconds())%len(sparkles)]
	sparkle2 := sparkles[(int(time.Since(kt.AnimationTime).Seconds())+3)%len(sparkles)]

	return fmt.Sprintf("%s %s %s", sparkle1, kt.Styles.Sparkle.Render(text), sparkle2)
}

// CreateFloatingEffect creates a floating box effect
func (kt *KawaiiTheme) CreateFloatingEffect(content string) string {
	return kt.Styles.FloatingBox.Render(content)
}

// CreateGlowEffect creates a glowing text effect
func (kt *KawaiiTheme) CreateGlowEffect(text string) string {
	return kt.Styles.Glow.Render(fmt.Sprintf("✨ %s ✨", text))
}

// GetWelcomeMessage returns a themed welcome message
func (kt *KawaiiTheme) GetWelcomeMessage() string {
	switch kt.Name {
	case "Sakura Dreams 🌸✨":
		return kt.CreateSparkleText("Welcome to your kawaii terminal paradise! 🌸")
	case "Galaxy Dreams 🌌⭐":
		return kt.CreateSparkleText("Welcome to the cosmic kawaii dimension! 🌌")
	case "Cyber Kawaii 🤖💫":
		return kt.CreateSparkleText("Initializing kawaii cybernetic interface! 🤖")
	case "Ocean Breeze 🌊🐚":
		return kt.CreateSparkleText("Dive into your oceanic kawaii world! 🌊")
	case "Rainbow Magic 🌈✨":
		return kt.ApplyRainbowEffect("Welcome to Rainbow Kawaii Land! ") + "🌈"
	default:
		return kt.CreateSparkleText("Welcome to Kawaii Shell! ✨")
	}
}

// Update updates theme animations
func (kt *KawaiiTheme) Update() {
	// Update animation time for gradient effects
	kt.AnimationTime = time.Now()
}

// GetThemes returns all available stunning themes
func GetThemes() []*KawaiiTheme {
	return []*KawaiiTheme{
		NewSakuraTheme(),
		NewGalaxyTheme(),
		NewCyberTheme(),
		NewOceanTheme(),
		NewRainbowTheme(),
	}
}
