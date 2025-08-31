package components

import (
	"math"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Particle represents a visual particle for effects
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Emoji    string
	Color    string
	Size     float64
	Rotation float64
}

// ParticleSystem manages particle effects
type ParticleSystem struct {
	particles []Particle
	width     int
	height    int
	active    bool
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(width, height int) *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0),
		width:     width,
		height:    height,
		active:    true,
	}
}

// SparkleEmoji returns random sparkle emojis
func SparkleEmoji() string {
	sparkles := []string{"✨", "⭐", "💫", "🌟", "✦", "✧", "⚡"}
	return sparkles[rand.Intn(len(sparkles))]
}

// HeartEmoji returns random heart emojis
func HeartEmoji() string {
	hearts := []string{"💕", "💖", "💗", "💓", "💝", "💘", "💞"}
	return hearts[rand.Intn(len(hearts))]
}

// FlowerEmoji returns random flower emojis
func FlowerEmoji() string {
	flowers := []string{"🌸", "🌺", "🌻", "🌷", "🌹", "🌼", "🌿"}
	return flowers[rand.Intn(len(flowers))]
}

// MagicEmoji returns random magic emojis
func MagicEmoji() string {
	magic := []string{"🔮", "🪄", "✨", "🌟", "⭐", "💫", "🎆", "🎇", "🌈", "🦄"}
	return magic[rand.Intn(len(magic))]
}

// FireworkEmoji returns random firework emojis
func FireworkEmoji() string {
	fireworks := []string{"🎆", "🎇", "✨", "💥", "🌟", "⚡", "💫"}
	return fireworks[rand.Intn(len(fireworks))]
}

// CelebrationEmoji returns random celebration emojis
func CelebrationEmoji() string {
	celebration := []string{"🎉", "🎊", "🥳", "🎈", "🎁", "🏆", "👑", "💎"}
	return celebration[rand.Intn(len(celebration))]
}

// ParticleType represents different types of particle effects
type ParticleType int

const (
	ParticleSparkle ParticleType = iota
	ParticleHeart
	ParticleFlower
	ParticleMagic
	ParticleFirework
	ParticleCelebration
	ParticleRainbow
	ParticleStardust
)

// AdvancedParticle represents a particle with advanced properties
type AdvancedParticle struct {
	Particle
	Type       ParticleType
	Trail      []Position
	Gravity    float64
	Bounce     float64
	Magnetism  float64
	Spiral     float64
	PulseRate  float64
	GlowRadius float64
	MagicPower float64
}

// Position represents a 2D position
type Position struct {
	X, Y float64
}

// AddSparkles adds sparkle particles around a point
func (ps *ParticleSystem) AddSparkles(x, y int, count int) {
	if !ps.active {
		return
	}

	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 0.5
		life := rand.Float64()*2 + 1

		particle := Particle{
			X:        float64(x) + rand.Float64()*4 - 2,
			Y:        float64(y) + rand.Float64()*4 - 2,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     life,
			MaxLife:  life,
			Emoji:    SparkleEmoji(),
			Size:     rand.Float64()*0.5 + 0.5,
			Rotation: rand.Float64() * 2 * math.Pi,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddHearts adds heart particles for happiness
func (ps *ParticleSystem) AddHearts(x, y int, count int) {
	if !ps.active {
		return
	}

	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*1.5 + 0.3
		life := rand.Float64()*3 + 2

		particle := Particle{
			X:       float64(x) + rand.Float64()*6 - 3,
			Y:       float64(y) + rand.Float64()*6 - 3,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle)*speed - 0.5, // Hearts float up
			Life:    life,
			MaxLife: life,
			Emoji:   HeartEmoji(),
			Size:    rand.Float64()*0.7 + 0.8,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddFlowerPetals adds flower petal effects
func (ps *ParticleSystem) AddFlowerPetals(x, y int, count int) {
	if !ps.active {
		return
	}

	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*1 + 0.2
		life := rand.Float64()*4 + 3

		particle := Particle{
			X:        float64(x) + rand.Float64()*8 - 4,
			Y:        float64(y) + rand.Float64()*8 - 4,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle)*speed*0.5 + 0.3, // Petals drift down
			Life:     life,
			MaxLife:  life,
			Emoji:    FlowerEmoji(),
			Size:     rand.Float64()*0.6 + 0.4,
			Rotation: rand.Float64() * 2 * math.Pi,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// Update updates all particles
func (ps *ParticleSystem) Update(deltaTime float64) {
	if !ps.active {
		return
	}

	// Update existing particles
	var alive []Particle
	for _, p := range ps.particles {
		// Update position
		p.X += p.VX * deltaTime
		p.Y += p.VY * deltaTime

		// Apply gravity and air resistance
		p.VY += 0.5 * deltaTime // Gravity
		p.VX *= 0.98            // Air resistance
		p.VY *= 0.98

		// Update life
		p.Life -= deltaTime

		// Update rotation
		p.Rotation += deltaTime * 2

		// Keep alive particles
		if p.Life > 0 && p.X >= 0 && p.X < float64(ps.width) && p.Y >= 0 && p.Y < float64(ps.height) {
			alive = append(alive, p)
		}
	}

	ps.particles = alive
}

// Render renders all particles to a string grid
func (ps *ParticleSystem) Render() [][]string {
	if !ps.active {
		return nil
	}

	// Create grid
	grid := make([][]string, ps.height)
	for i := range grid {
		grid[i] = make([]string, ps.width)
	}

	// Render particles
	for _, p := range ps.particles {
		x := int(math.Round(p.X))
		y := int(math.Round(p.Y))

		if x >= 0 && x < ps.width && y >= 0 && y < ps.height {
			// Apply alpha based on life
			alpha := p.Life / p.MaxLife
			if alpha > 0.1 { // Only show if visible enough
				grid[y][x] = p.Emoji
			}
		}
	}

	return grid
}

// Clear clears all particles
func (ps *ParticleSystem) Clear() {
	ps.particles = ps.particles[:0]
}

// SetActive enables/disables the particle system
func (ps *ParticleSystem) SetActive(active bool) {
	ps.active = active
}

// ParticleTickMsg represents a particle update tick
type ParticleTickMsg struct {
	Time time.Time
}

// ParticleUpdateCmd returns a command for particle updates
func ParticleUpdateCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return ParticleTickMsg{Time: t}
	})
}

// AddMagicBlast creates a magical explosion effect
func (ps *ParticleSystem) AddMagicBlast(x, y int, intensity int) {
	if !ps.active {
		return
	}

	count := 20 + intensity*5
	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*4 + 2
		life := rand.Float64()*3 + 2

		particle := Particle{
			X:        float64(x),
			Y:        float64(y),
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     life,
			MaxLife:  life,
			Emoji:    MagicEmoji(),
			Size:     rand.Float64()*0.8 + 0.7,
			Rotation: rand.Float64() * 2 * math.Pi,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddFireworks creates a firework explosion
func (ps *ParticleSystem) AddFireworks(x, y int, colors []string) {
	if !ps.active {
		return
	}

	// Main burst
	for i := 0; i < 25; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*3 + 1.5
		life := rand.Float64()*4 + 3

		particle := Particle{
			X:       float64(x),
			Y:       float64(y),
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Life:    life,
			MaxLife: life,
			Emoji:   FireworkEmoji(),
			Size:    rand.Float64()*1.2 + 0.8,
		}

		if len(colors) > 0 {
			particle.Color = colors[rand.Intn(len(colors))]
		}

		ps.particles = append(ps.particles, particle)
	}

	// Secondary sparkles
	for i := 0; i < 15; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*1.5 + 0.5
		life := rand.Float64()*2 + 1.5

		particle := Particle{
			X:       float64(x) + rand.Float64()*10 - 5,
			Y:       float64(y) + rand.Float64()*10 - 5,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Life:    life,
			MaxLife: life,
			Emoji:   SparkleEmoji(),
			Size:    rand.Float64()*0.6 + 0.4,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddRainbowTrail creates a rainbow particle trail
func (ps *ParticleSystem) AddRainbowTrail(startX, startY, endX, endY int) {
	if !ps.active {
		return
	}

	colors := []string{
		"#ff0000", // Red
		"#ff8000", // Orange
		"#ffff00", // Yellow
		"#00ff00", // Green
		"#0080ff", // Blue
		"#8000ff", // Purple
		"#ff00ff", // Magenta
	}

	deltaX := float64(endX - startX)
	deltaY := float64(endY - startY)
	steps := int(math.Sqrt(deltaX*deltaX + deltaY*deltaY))

	for i := 0; i <= steps; i++ {
		progress := float64(i) / float64(steps)
		x := float64(startX) + deltaX*progress
		y := float64(startY) + deltaY*progress

		// Add multiple particles at each step
		for j := 0; j < 3; j++ {
			life := rand.Float64()*2 + 1
			particle := Particle{
				X:       x + rand.Float64()*4 - 2,
				Y:       y + rand.Float64()*4 - 2,
				VX:      rand.Float64()*0.5 - 0.25,
				VY:      rand.Float64()*0.5 - 0.25,
				Life:    life,
				MaxLife: life,
				Emoji:   "✨",
				Color:   colors[i%len(colors)],
				Size:    rand.Float64()*0.7 + 0.3,
			}

			ps.particles = append(ps.particles, particle)
		}
	}
}

// AddStardust creates a magical stardust effect
func (ps *ParticleSystem) AddStardust(x, y, width, height int, density int) {
	if !ps.active {
		return
	}

	for i := 0; i < density; i++ {
		px := float64(x) + rand.Float64()*float64(width)
		py := float64(y) + rand.Float64()*float64(height)
		life := rand.Float64()*5 + 3

		particle := Particle{
			X:        px,
			Y:        py,
			VX:       rand.Float64()*0.3 - 0.15,
			VY:       -rand.Float64()*0.5 - 0.2, // Gentle upward drift
			Life:     life,
			MaxLife:  life,
			Emoji:    SparkleEmoji(),
			Size:     rand.Float64()*0.5 + 0.3,
			Rotation: rand.Float64() * 2 * math.Pi,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddCelebrationBurst creates a massive celebration effect
func (ps *ParticleSystem) AddCelebrationBurst(x, y int) {
	if !ps.active {
		return
	}

	// Central explosion
	ps.AddFireworks(x, y, []string{
		"#ffd700",
		"#ff6b6b",
		"#4ecdc4",
		"#45b7d1",
		"#96ceb4",
	})

	// Surrounding celebrations
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		offsetX := int(math.Cos(angle) * 8)
		offsetY := int(math.Sin(angle) * 8)

		for j := 0; j < 10; j++ {
			life := rand.Float64()*3 + 2
			speed := rand.Float64()*2 + 0.5

			particle := Particle{
				X:       float64(x + offsetX),
				Y:       float64(y + offsetY),
				VX:      math.Cos(angle+rand.Float64()*0.5-0.25) * speed,
				VY:      math.Sin(angle+rand.Float64()*0.5-0.25) * speed,
				Life:    life,
				MaxLife: life,
				Emoji:   CelebrationEmoji(),
				Size:    rand.Float64()*1.0 + 0.5,
			}

			ps.particles = append(ps.particles, particle)
		}
	}

	// Add hearts and sparkles
	ps.AddHearts(x, y, 15)
	ps.AddSparkles(x, y, 25)
}

// AddSpiralEffect creates a magical spiral effect
func (ps *ParticleSystem) AddSpiralEffect(centerX, centerY int, radius float64, turns int) {
	if !ps.active {
		return
	}

	steps := turns * 20
	angleStep := float64(turns) * 2 * math.Pi / float64(steps)
	radiusStep := radius / float64(steps)

	for i := 0; i < steps; i++ {
		angle := float64(i) * angleStep
		currentRadius := float64(i) * radiusStep

		x := float64(centerX) + math.Cos(angle)*currentRadius
		y := float64(centerY) + math.Sin(angle)*currentRadius

		life := rand.Float64()*2 + 1.5
		particle := Particle{
			X:        x,
			Y:        y,
			VX:       math.Cos(angle+math.Pi/2) * 0.5, // Perpendicular to radius
			VY:       math.Sin(angle+math.Pi/2) * 0.5,
			Life:     life,
			MaxLife:  life,
			Emoji:    MagicEmoji(),
			Size:     rand.Float64()*0.6 + 0.4,
			Rotation: angle,
		}

		ps.particles = append(ps.particles, particle)
	}
}

// AddWaveEffect creates a wave propagation effect
func (ps *ParticleSystem) AddWaveEffect(centerX, centerY, maxRadius int, intensity int) {
	if !ps.active {
		return
	}

	waves := 5
	for wave := 0; wave < waves; wave++ {
		waveRadius := float64(wave+1) * float64(maxRadius) / float64(waves)
		particlesInWave := 15 + intensity*3

		for i := 0; i < particlesInWave; i++ {
			angle := float64(i) * 2 * math.Pi / float64(particlesInWave)
			x := float64(centerX) + math.Cos(angle)*waveRadius
			y := float64(centerY) + math.Sin(angle)*waveRadius

			life := rand.Float64()*2 + 1 + float64(wave)*0.3
			particle := Particle{
				X:       x,
				Y:       y,
				VX:      math.Cos(angle) * 0.3,
				VY:      math.Sin(angle) * 0.3,
				Life:    life,
				MaxLife: life,
				Emoji:   SparkleEmoji(),
				Size:    rand.Float64()*0.5 + 0.3,
			}

			ps.particles = append(ps.particles, particle)
		}
	}
}

// AddMagicCircle creates a magical circle effect
func (ps *ParticleSystem) AddMagicCircle(centerX, centerY, radius int, rotationSpeed float64) {
	if !ps.active {
		return
	}

	particleCount := 16
	for i := 0; i < particleCount; i++ {
		baseAngle := float64(i) * 2 * math.Pi / float64(particleCount)

		x := float64(centerX) + math.Cos(baseAngle)*float64(radius)
		y := float64(centerY) + math.Sin(baseAngle)*float64(radius)

		life := rand.Float64()*4 + 3
		particle := Particle{
			X:        x,
			Y:        y,
			VX:       math.Cos(baseAngle+math.Pi/2) * rotationSpeed,
			VY:       math.Sin(baseAngle+math.Pi/2) * rotationSpeed,
			Life:     life,
			MaxLife:  life,
			Emoji:    MagicEmoji(),
			Size:     rand.Float64()*0.8 + 0.6,
			Rotation: baseAngle,
		}

		ps.particles = append(ps.particles, particle)
	}

	// Add central sparkle
	life := rand.Float64()*3 + 2
	centerParticle := Particle{
		X:       float64(centerX),
		Y:       float64(centerY),
		VX:      0,
		VY:      0,
		Life:    life,
		MaxLife: life,
		Emoji:   "🔮",
		Size:    1.5,
	}

	ps.particles = append(ps.particles, centerParticle)
}

// CreateMagicalAura creates a continuous magical aura around a point
func (ps *ParticleSystem) CreateMagicalAura(centerX, centerY, radius int) tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(time.Time) tea.Msg {
		if ps.active {
			// Add gentle sparkles in a circle
			angle := rand.Float64() * 2 * math.Pi
			distance := rand.Float64() * float64(radius)
			x := int(float64(centerX) + math.Cos(angle)*distance)
			y := int(float64(centerY) + math.Sin(angle)*distance)

			ps.AddSparkles(x, y, 1)
		}
		return ParticleTickMsg{Time: time.Now()}
	})
}
