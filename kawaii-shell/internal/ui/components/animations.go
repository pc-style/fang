package components

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

type AnimationState int

const (
	AnimIdle AnimationState = iota
	AnimBounce
	AnimWiggle
	AnimSpin
	AnimFloat
	AnimPulse
	AnimGlow
)

type AnimatedElement struct {
	Content       string
	X, Y          float64
	TargetX       float64
	TargetY       float64
	Scale         float64
	Rotation      float64
	Alpha         float64
	State         AnimationState
	Time          float64
	Duration      float64
	Easing        EasingFunc
	Style         lipgloss.Style
	OnComplete    func()
	BounceHeight  float64
	GlowIntensity float64
}

type EasingFunc func(t float64) float64

func EaseInOut(t float64) float64 {
	return t * t * (3 - 2*t)
}

func EaseBounce(t float64) float64 {
	if t < 1/2.75 {
		return 7.5625 * t * t
	} else if t < 2/2.75 {
		t -= 1.5 / 2.75
		return 7.5625*t*t + 0.75
	} else if t < 2.5/2.75 {
		t -= 2.25 / 2.75
		return 7.5625*t*t + 0.9375
	} else {
		t -= 2.625 / 2.75
		return 7.5625*t*t + 0.984375
	}
}

func EaseElastic(t float64) float64 {
	return math.Pow(2, -10*t)*math.Sin((t-0.1)*5*math.Pi) + 1
}

func EaseSine(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

func NewAnimatedElement(content string, x, y float64) *AnimatedElement {
	return &AnimatedElement{
		Content:  content,
		X:        x,
		Y:        y,
		TargetX:  x,
		TargetY:  y,
		Scale:    1.0,
		Rotation: 0,
		Alpha:    1.0,
		State:    AnimIdle,
		Time:     0,
		Duration: 1.0,
		Easing:   EaseInOut,
		Style:    lipgloss.NewStyle(),
	}
}

func (ae *AnimatedElement) MoveTo(targetX, targetY float64, duration float64) *AnimatedElement {
	ae.TargetX = targetX
	ae.TargetY = targetY
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimFloat
	ae.Easing = EaseInOut
	return ae
}

func (ae *AnimatedElement) Bounce(height float64, duration float64) *AnimatedElement {
	ae.BounceHeight = height
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimBounce
	ae.Easing = EaseBounce
	return ae
}

func (ae *AnimatedElement) Wiggle(duration float64) *AnimatedElement {
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimWiggle
	ae.Easing = EaseSine
	return ae
}

func (ae *AnimatedElement) Spin(duration float64) *AnimatedElement {
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimSpin
	ae.Easing = EaseInOut
	return ae
}

func (ae *AnimatedElement) Pulse(duration float64) *AnimatedElement {
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimPulse
	ae.Easing = EaseSine
	return ae
}

func (ae *AnimatedElement) Glow(intensity float64, duration float64) *AnimatedElement {
	ae.GlowIntensity = intensity
	ae.Duration = duration
	ae.Time = 0
	ae.State = AnimGlow
	ae.Easing = EaseSine
	return ae
}

func (ae *AnimatedElement) Update(deltaTime float64) {
	if ae.State == AnimIdle {
		return
	}
	ae.Time += deltaTime
	progress := ae.Time / ae.Duration
	if progress >= 1.0 {
		progress = 1.0
		ae.State = AnimIdle
		if ae.OnComplete != nil {
			ae.OnComplete()
		}
	}
	easedProgress := ae.Easing(progress)
	switch ae.State {
	case AnimFloat:
		ae.X = ae.X + (ae.TargetX-ae.X)*easedProgress*0.1
		ae.Y = ae.Y + (ae.TargetY-ae.Y)*easedProgress*0.1
	case AnimBounce:
		ae.Y = ae.TargetY - ae.BounceHeight*easedProgress
	case AnimWiggle:
		ae.X = ae.TargetX + math.Sin(progress*math.Pi*4)*5*easedProgress
	case AnimSpin:
		ae.Rotation = progress * 2 * math.Pi
	case AnimPulse:
		ae.Scale = 1.0 + math.Sin(progress*math.Pi*2)*0.3*easedProgress
	case AnimGlow:
		ae.Alpha = 1.0 + math.Sin(progress*math.Pi*2)*ae.GlowIntensity*easedProgress
		if ae.Alpha > 1.0 {
			ae.Alpha = 1.0
		}
		if ae.Alpha < 0.1 {
			ae.Alpha = 0.1
		}
	}
}

func (ae *AnimatedElement) Render() string {
	style := ae.Style
	if ae.Alpha < 1.0 {
		style = style.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", int(255*ae.Alpha), int(255*ae.Alpha), int(255*ae.Alpha))))
	}
	if ae.Scale != 1.0 {
		padding := int((ae.Scale - 1.0) * 2)
		if padding > 0 {
			style = style.Padding(padding, padding)
		}
	}
	content := ae.Content
	if ae.State == AnimGlow && ae.Alpha > 1.0 {
		content = fmt.Sprintf("✨%s✨", content)
	}
	if ae.State == AnimBounce {
		content = fmt.Sprintf("⬆️%s⬆️", content)
	}
	if ae.State == AnimWiggle {
		content = fmt.Sprintf("🐍%s🐍", content)
	}
	if ae.State == AnimSpin {
		content = fmt.Sprintf("🌀%s🌀", content)
	}
	return style.Render(content)
}

type AnimationManager struct {
	elements   []*AnimatedElement
	lastUpdate time.Time
	particles  *ParticleSystem
}

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		elements:   make([]*AnimatedElement, 0),
		lastUpdate: time.Now(),
	}
}

func (am *AnimationManager) AddElement(element *AnimatedElement) {
	am.elements = append(am.elements, element)
}

func (am *AnimationManager) Update() {
	now := time.Now()
	deltaTime := now.Sub(am.lastUpdate).Seconds()
	am.lastUpdate = now
	for _, element := range am.elements {
		element.Update(deltaTime)
	}
	if am.particles != nil {
		am.particles.Update(deltaTime)
	}
}

func (am *AnimationManager) Render() string {
	var parts []string
	for _, element := range am.elements {
		parts = append(parts, element.Render())
	}
	return strings.Join(parts, " ")
}

func (am *AnimationManager) Clear() {
	am.elements = am.elements[:0]
}

func (am *AnimationManager) SetParticleSystem(ps *ParticleSystem) {
	am.particles = ps
}

func CreateSparkleText(text string, style lipgloss.Style) *AnimatedElement {
	element := NewAnimatedElement(text, 0, 0)
	element.Style = style
	element.Glow(0.5, 2.0)
	return element
}

func CreateBouncyButton(text string, style lipgloss.Style) *AnimatedElement {
	element := NewAnimatedElement(text, 0, 0)
	element.Style = style.Border(lipgloss.RoundedBorder()).Padding(0, 1)
	element.Bounce(2, 0.5)
	return element
}

func CreateFloatingText(text string, style lipgloss.Style) *AnimatedElement {
	element := NewAnimatedElement(text, 0, 0)
	element.Style = style
	element.MoveTo(0, -5, 2.0)
	element.Easing = EaseElastic
	return element
}

type TransitionType int

const (
	TransitionCrossfade TransitionType = iota
	TransitionSlide
	TransitionMorph
	TransitionScramble
)

type Transition struct {
	From      string
	To        string
	Style     lipgloss.Style
	Type      TransitionType
	Duration  float64
	Easing    EasingFunc
	Started   time.Time
	Progress  float64
	Completed bool
}

func NewTransition(from, to string, style lipgloss.Style, typ TransitionType, duration float64) *Transition {
	return &Transition{
		From:     from,
		To:       to,
		Style:    style,
		Type:     typ,
		Duration: duration,
		Easing:   EaseInOut,
		Started:  time.Now(),
	}
}

func (t *Transition) Update(delta float64) {
	if t.Completed {
		return
	}
	t.Progress += delta / t.Duration
	if t.Progress >= 1.0 {
		t.Progress = 1.0
		t.Completed = true
	}
}

func (t *Transition) Render() string {
	p := t.Easing(t.Progress)
	switch t.Type {
	case TransitionCrossfade:
		return t.renderCrossfade(p)
	case TransitionSlide:
		return t.renderSlide(p)
	case TransitionMorph:
		return t.renderMorph(p)
	case TransitionScramble:
		return t.renderScramble(p)
	default:
		return t.Style.Render(t.To)
	}
}

func (t *Transition) renderCrossfade(p float64) string {
	from := []rune(t.From)
	to := []rune(t.To)
	max := len(to)
	if len(from) > max {
		max = len(from)
	}
	cut := int(p * float64(max))
	var b strings.Builder
	for i := 0; i < max; i++ {
		var r rune
		if i < cut {
			if i < len(to) {
				r = to[i]
			} else {
				r = ' '
			}
		} else {
			if i < len(from) {
				r = from[i]
			} else {
				r = ' '
			}
		}
		b.WriteRune(r)
	}
	return t.Style.Render(b.String())
}

func (t *Transition) renderSlide(p float64) string {
	width := len([]rune(t.To))
	offset := int((1 - p) * float64(width))
	if offset < 0 {
		offset = 0
	}
	padding := strings.Repeat(" ", offset)
	return t.Style.Render(padding + t.To)
}

func (t *Transition) renderMorph(p float64) string {
	from := []rune(t.From)
	to := []rune(t.To)
	max := len(to)
	if len(from) > max {
		max = len(from)
	}
	threshold := int(p * float64(max))
	var b strings.Builder
	for i := 0; i < max; i++ {
		if i < threshold {
			if i < len(to) {
				b.WriteRune(to[i])
			} else {
				b.WriteRune(' ')
			}
		} else {
			if i < len(from) {
				b.WriteRune(from[i])
			} else {
				b.WriteRune(' ')
			}
		}
	}
	return t.Style.Render(b.String())
}

func (t *Transition) renderScramble(p float64) string {
	to := []rune(t.To)
	count := int((1 - p) * float64(len(to)))
	var b strings.Builder
	for i := 0; i < len(to); i++ {
		if i < count {
			b.WriteRune(randomRune())
		} else {
			b.WriteRune(to[i])
		}
	}
	return t.Style.Render(b.String())
}

func randomRune() rune {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()☼★✦✧❤✿")
	return alphabet[rand.Intn(len(alphabet))]
}
