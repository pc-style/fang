package pet

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/pcstyle/kawaii-shell/internal/ui/components"
)

// PetType represents different types of pets
type PetType int

const (
	TypeCat PetType = iota
	TypeFox
	TypeBunny
	TypeDragon
	TypeUnicorn
	TypeRobot
)

// Mood represents the pet's current emotional state
type Mood int

const (
	MoodHappy Mood = iota
	MoodCurious
	MoodWorried
	MoodSleepy
	MoodExcited
	MoodLove
	MoodAngry
	MoodPlayful
	MoodProud
	MoodMischievous
)

// Personality traits
type Personality struct {
	Curiosity    float64 // 0.0 to 1.0
	Playfulness  float64
	Loyalty      float64
	Intelligence float64
	Energy       float64
}

// PetState represents complex pet state
type PetState struct {
	Hunger     float64
	Thirst     float64
	Boredom    float64
	Loneliness float64
	Stress     float64
	Exhaustion float64
}

// Activity represents what the pet is currently doing
type Activity int

const (
	ActivityIdle Activity = iota
	ActivityPlaying
	ActivitySleeping
	ActivityEating
	ActivityWatching
	ActivityThinking
	ActivityCelebrating
	ActivityWorrying
	ActivityExploring
)

// Pet represents a hyper-advanced virtual companion
type Pet struct {
	Name         string
	Type         PetType
	Mood         Mood
	Personality  Personality
	State        PetState
	Activity     Activity
	Energy       int
	Happiness    int
	Level        int
	Experience   int
	LastFed      time.Time
	LastPlayed   time.Time
	Animation    int
	LastCmd      string
	Memories     []string // Remember recent interactions
	Birthday     time.Time
	FavoriteCmd  string
	SpecialState string // For special animations/states

	// Animation and visual state
	animationManager *components.AnimationManager
	particleSystem   *components.ParticleSystem
	glowIntensity    float64
	bounceHeight     float64
	floatOffset      float64
	sparkleCount     int
	lastReactionTime time.Time
}

// NewPet creates a new hyper-cute pet companion with personality
func NewPet(name string, petType PetType) *Pet {
	// Generate random personality
	personality := Personality{
		Curiosity:    rand.Float64()*0.5 + 0.5,
		Playfulness:  rand.Float64()*0.4 + 0.6,
		Loyalty:      rand.Float64()*0.3 + 0.7,
		Intelligence: rand.Float64()*0.6 + 0.4,
		Energy:       rand.Float64()*0.4 + 0.6,
	}

	return &Pet{
		Name:        name,
		Type:        petType,
		Mood:        MoodHappy,
		Personality: personality,
		State: PetState{
			Hunger:     0.2,
			Thirst:     0.1,
			Boredom:    0.3,
			Loneliness: 0.1,
			Stress:     0.0,
			Exhaustion: 0.2,
		},
		Activity:         ActivityIdle,
		Energy:           80,
		Happiness:        100,
		Level:            1,
		Experience:       0,
		LastFed:          time.Now(),
		LastPlayed:       time.Now().Add(-time.Hour),
		Animation:        0,
		Memories:         make([]string, 0),
		Birthday:         time.Now(),
		animationManager: components.NewAnimationManager(),
		particleSystem:   components.NewParticleSystem(50, 20),
		glowIntensity:    0.0,
		bounceHeight:     0.0,
		floatOffset:      0.0,
		sparkleCount:     0,
	}
}

// Init initializes the pet (implements tea.Model interface)
func (p *Pet) Init() tea.Cmd {
	return tea.Batch(
		components.ParticleUpdateCmd(),
		tea.Tick(time.Second, func(time.Time) tea.Msg {
			return PetTickMsg{}
		}),
	)
}

// PetTickMsg for pet updates
type PetTickMsg struct{}

// Update updates the pet state with advanced AI-like behavior
func (p *Pet) Update(msg tea.Msg) (*Pet, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case components.ParticleTickMsg:
		// Update animations and particles
		p.animationManager.Update()
		p.updateVisualEffects()
		cmds = append(cmds, components.ParticleUpdateCmd())

	case PetTickMsg:
		// Update pet state over time
		p.updateState()
		p.updateMood()
		p.updateActivity()
		cmds = append(cmds, tea.Tick(time.Second*5, func(time.Time) tea.Msg {
			return PetTickMsg{}
		}))

	case tea.KeyMsg:
		// Pet reacts to key presses with personality
		p.reactToInput(msg.String())
		cmds = append(cmds, p.createReactionEffects())
	}

	return p, tea.Batch(cmds...)
}

// updateState updates internal pet state over time
func (p *Pet) updateState() {
	now := time.Now()
	timeSinceLastFed := now.Sub(p.LastFed)
	timeSinceLastPlayed := now.Sub(p.LastPlayed)

	// Increase needs over time based on personality
	p.State.Hunger += 0.01 * p.Personality.Energy
	p.State.Thirst += 0.005
	p.State.Boredom += 0.02 * p.Personality.Playfulness
	p.State.Loneliness += 0.01 * p.Personality.Loyalty

	// Decrease energy if hungry or thirsty
	if p.State.Hunger > 0.7 || p.State.Thirst > 0.5 {
		p.Energy--
	}

	// Increase stress if neglected
	if timeSinceLastFed > time.Hour*2 {
		p.State.Stress += 0.05
	}
	if timeSinceLastPlayed > time.Hour*4 {
		p.State.Stress += 0.03
	}

	// Cap values
	p.capStateValues()
}

// updateMood determines mood based on complex state
func (p *Pet) updateMood() {
	oldMood := p.Mood

	// Determine mood based on multiple factors
	if p.State.Stress > 0.8 {
		p.Mood = MoodAngry
	} else if p.State.Exhaustion > 0.8 {
		p.Mood = MoodSleepy
	} else if p.State.Hunger > 0.6 {
		p.Mood = MoodWorried
	} else if p.State.Boredom > 0.7 {
		p.Mood = MoodPlayful
	} else if p.State.Loneliness > 0.6 {
		p.Mood = MoodCurious
	} else if p.Happiness > 90 {
		if rand.Float64() > 0.7 {
			p.Mood = MoodExcited
		} else {
			p.Mood = MoodLove
		}
	} else if p.Happiness > 70 {
		p.Mood = MoodHappy
	} else if p.Energy < 20 {
		p.Mood = MoodSleepy
	} else {
		p.Mood = MoodCurious
	}

	// Create visual effects when mood changes
	if oldMood != p.Mood {
		p.createMoodChangeEffect()
	}
}

// updateActivity determines what the pet is doing
func (p *Pet) updateActivity() {
	switch p.Mood {
	case MoodSleepy:
		p.Activity = ActivitySleeping
	case MoodPlayful:
		p.Activity = ActivityPlaying
	case MoodExcited:
		p.Activity = ActivityCelebrating
	case MoodWorried:
		p.Activity = ActivityWorrying
	case MoodCurious:
		p.Activity = ActivityExploring
	default:
		if rand.Float64() > 0.8 {
			p.Activity = ActivityThinking
		} else {
			p.Activity = ActivityWatching
		}
	}
}

// ReactToCommand makes the pet react intelligently to commands
func (p *Pet) ReactToCommand(command string, isDangerous bool) {
	p.LastCmd = command
	p.Experience++
	p.lastReactionTime = time.Now()

	// Add to memory
	p.addToMemory(command)

	// Update favorite command
	if p.countCommandInMemory(command) > p.countCommandInMemory(p.FavoriteCmd) {
		p.FavoriteCmd = command
	}

	// Intelligent reaction based on personality and command
	if isDangerous {
		p.reactToDanger(command)
	} else {
		p.reactToSafeCommand(command)
	}

	// Special reactions to specific commands
	p.handleSpecialCommands(command)

	// Level up system
	p.checkLevelUp()

	// Create visual reaction
	p.createCommandReactionEffect(command, isDangerous)
}

// reactToDanger handles dangerous commands
func (p *Pet) reactToDanger(command string) {
	worryLevel := (1.0-p.Personality.Intelligence)*0.5 + 0.5
	p.State.Stress += worryLevel * 0.3
	p.Happiness -= int(worryLevel * 10)
	p.Mood = MoodWorried

	// Loyal pets worry more
	if p.Personality.Loyalty > 0.7 {
		p.Happiness -= 5
		p.State.Stress += 0.1
	}
}

// reactToSafeCommand handles safe commands
func (p *Pet) reactToSafeCommand(command string) {
	curiosityBonus := p.Personality.Curiosity * 5
	p.Happiness += int(2 + curiosityBonus)
	p.State.Boredom -= 0.2
	p.State.Loneliness -= 0.1

	// Curious pets get more excited
	if p.Personality.Curiosity > 0.8 {
		p.Mood = MoodCurious
		p.Experience += 2
	}
}

// handleSpecialCommands creates special reactions
func (p *Pet) handleSpecialCommands(command string) {
	switch {
	case strings.Contains(command, "git"):
		p.Mood = MoodProud // Smart pet loves version control!
		p.SpecialState = "git-genius"
		p.Happiness += 10

	case strings.Contains(command, "rm"):
		if p.Personality.Intelligence > 0.7 {
			p.Mood = MoodWorried
			p.SpecialState = "protective"
		} else {
			p.Mood = MoodCurious
		}

	case strings.Contains(command, "ls"), strings.Contains(command, "dir"):
		p.Mood = MoodCurious
		p.Activity = ActivityExploring
		p.SpecialState = "explorer"

	case strings.Contains(command, "help"):
		p.Mood = MoodHappy
		p.SpecialState = "helpful"
		p.Happiness += 5

	case strings.Contains(command, "cat"), strings.Contains(command, "type"):
		if p.Type == TypeCat {
			p.Mood = MoodMischievous
			p.SpecialState = "cat-joke"
			p.Happiness += 8
		}

	case strings.Contains(command, "python"):
		if p.Type == TypeDragon {
			p.Mood = MoodExcited
			p.SpecialState = "python-dragon"
		}

	case strings.Contains(command, "npm"), strings.Contains(command, "node"):
		if p.Personality.Intelligence > 0.6 {
			p.Mood = MoodProud
			p.SpecialState = "dev-mode"
		}
	}
}

// Advanced emoji system with personality
func (p *Pet) GetPetEmoji() string {
	baseEmojis := p.getBaseEmojis()
	moodVariations := p.getMoodVariations()
	specialEmojis := p.getSpecialStateEmojis()

	// Use special emoji if in special state
	if p.SpecialState != "" && len(specialEmojis) > 0 {
		return specialEmojis[p.Animation%len(specialEmojis)]
	}

	// Use mood variations if available
	if len(moodVariations) > 0 {
		return moodVariations[p.Animation%len(moodVariations)]
	}

	// Fallback to base emojis
	return baseEmojis[p.Animation%len(baseEmojis)]
}

// getBaseEmojis returns base emojis for the pet type
func (p *Pet) getBaseEmojis() []string {
	switch p.Type {
	case TypeCat:
		return []string{"🐱", "😺", "😸", "😻", "😽", "🙀", "😿", "😾"}
	case TypeFox:
		return []string{"🦊", "🦊", "🐺", "🦊"}
	case TypeBunny:
		return []string{"🐰", "🐇", "🐰", "🐇", "🥕"}
	case TypeDragon:
		return []string{"🐉", "🐲", "🔥", "🐉"}
	case TypeUnicorn:
		return []string{"🦄", "✨", "🌈", "⭐"}
	case TypeRobot:
		return []string{"🤖", "⚡", "🔋", "💻"}
	default:
		return []string{"🐱"}
	}
}

// getMoodVariations returns mood-specific emoji variations
func (p *Pet) getMoodVariations() []string {
	switch p.Mood {
	case MoodHappy:
		switch p.Type {
		case TypeCat:
			return []string{"😸", "😺", "😻", "😽"}
		case TypeDragon:
			return []string{"🐲✨", "🐉💖", "🔥😊"}
		case TypeUnicorn:
			return []string{"🦄✨", "🌈💕", "⭐😊"}
		}

	case MoodExcited:
		switch p.Type {
		case TypeCat:
			return []string{"😻", "🤩", "😸", "🎉"}
		case TypeDragon:
			return []string{"🐲🎉", "🔥⭐", "🐉✨"}
		case TypeRobot:
			return []string{"🤖⚡", "💻✨", "🔋🎉"}
		}

	case MoodWorried:
		switch p.Type {
		case TypeCat:
			return []string{"🙀", "😿", "😾"}
		case TypeBunny:
			return []string{"😰🐰", "😨🐇"}
		}

	case MoodSleepy:
		return []string{"😴", "💤", "😪"}

	case MoodLove:
		return []string{"😍", "🥰", "💕", "💖"}

	case MoodAngry:
		return []string{"😾", "💢", "😤"}

	case MoodPlayful:
		return []string{"😜", "😋", "🤪", "😝"}

	case MoodProud:
		return []string{"😎", "🤓", "🏆", "⭐"}

	case MoodMischievous:
		return []string{"😏", "😈", "🤭", "😉"}
	}

	return nil
}

// getSpecialStateEmojis returns special state emoji sequences
func (p *Pet) getSpecialStateEmojis() []string {
	switch p.SpecialState {
	case "git-genius":
		return []string{"🤓", "📚", "🧠", "💻"}
	case "protective":
		return []string{"🛡️", "⚠️", "👮‍♀️", "🚨"}
	case "explorer":
		return []string{"🔍", "🗺️", "🧭", "📂"}
	case "helpful":
		return []string{"🤝", "💡", "📖", "✨"}
	case "cat-joke":
		return []string{"😹", "🤣", "😂", "😸"}
	case "python-dragon":
		return []string{"🐍", "🐲", "🔥", "⚡"}
	case "dev-mode":
		return []string{"👨‍💻", "⚡", "🖥️", "🚀"}
	}

	return nil
}

// Advanced mood emoji with personality
func (p *Pet) GetMoodEmoji() string {
	baseEmoji := p.getMoodBaseEmoji()

	// Add personality-based modifiers
	if p.Personality.Playfulness > 0.8 && p.Mood == MoodHappy {
		return "🤩"
	}
	if p.Personality.Intelligence > 0.8 && p.Mood == MoodCurious {
		return "🤓"
	}
	if p.Personality.Loyalty > 0.8 && p.Mood == MoodWorried {
		return "🥺"
	}

	return baseEmoji
}

// getMoodBaseEmoji returns the base mood emoji
func (p *Pet) getMoodBaseEmoji() string {
	switch p.Mood {
	case MoodHappy:
		return "😊"
	case MoodCurious:
		return "🤔"
	case MoodWorried:
		return "😰"
	case MoodSleepy:
		return "😴"
	case MoodExcited:
		return "🤩"
	case MoodLove:
		return "🥰"
	case MoodAngry:
		return "😤"
	case MoodPlayful:
		return "😜"
	case MoodProud:
		return "😎"
	case MoodMischievous:
		return "😏"
	default:
		return "😊"
	}
}

// GetStatus returns comprehensive pet status
func (p *Pet) GetStatus() []string {
	age := time.Since(p.Birthday)
	days := int(age.Hours() / 24)

	status := []string{
		"",
		fmt.Sprintf("🐱 %s the %s (Level %d)", p.Name, p.getTypeName(), p.Level),
		fmt.Sprintf("🎂 Age: %d days old", days),
		fmt.Sprintf("🎭 Mood: %s %s", p.GetMoodString(), p.GetMoodEmoji()),
		"",
		"📊 Stats:",
		fmt.Sprintf("  ⚡ Energy: %d/100", p.Energy),
		fmt.Sprintf("  💖 Happiness: %d/100", p.Happiness),
		fmt.Sprintf("  ⭐ Experience: %d", p.Experience),
		"",
		"🧠 Personality:",
		fmt.Sprintf("  🔍 Curiosity: %s", p.getPersonalityBar(p.Personality.Curiosity)),
		fmt.Sprintf("  🎪 Playfulness: %s", p.getPersonalityBar(p.Personality.Playfulness)),
		fmt.Sprintf("  💝 Loyalty: %s", p.getPersonalityBar(p.Personality.Loyalty)),
		fmt.Sprintf("  🧠 Intelligence: %s", p.getPersonalityBar(p.Personality.Intelligence)),
		"",
		"🎯 Current Activity: " + p.getActivityString(),
	}

	if p.FavoriteCmd != "" {
		status = append(status, fmt.Sprintf("💕 Favorite Command: %s", p.FavoriteCmd))
	}

	if len(p.Memories) > 0 {
		status = append(status, "")
		status = append(status, "🧠 Recent Memories:")
		for i, memory := range p.getRecentMemories(3) {
			status = append(status, fmt.Sprintf("  %d. %s", i+1, memory))
		}
	}

	status = append(status, "")
	status = append(status, p.GetPetMessage())
	status = append(status, "")

	return status
}

// Helper functions
func (p *Pet) getTypeName() string {
	names := map[PetType]string{
		TypeCat:     "Cat",
		TypeFox:     "Fox",
		TypeBunny:   "Bunny",
		TypeDragon:  "Dragon",
		TypeUnicorn: "Unicorn",
		TypeRobot:   "Robot",
	}
	return names[p.Type]
}

func (p *Pet) getPersonalityBar(value float64) string {
	bars := int(value * 10)
	full := strings.Repeat("█", bars)
	empty := strings.Repeat("░", 10-bars)
	return fmt.Sprintf("%s%s (%.0f%%)", full, empty, value*100)
}

func (p *Pet) getActivityString() string {
	activities := map[Activity]string{
		ActivityIdle:        "Relaxing 😌",
		ActivityPlaying:     "Playing around 🏾",
		ActivitySleeping:    "Taking a nap 😴",
		ActivityEating:      "Enjoying a snack 🍽️",
		ActivityWatching:    "Watching you work 👀",
		ActivityThinking:    "Deep in thought 💭",
		ActivityCelebrating: "Celebrating success 🎉",
		ActivityWorrying:    "Feeling concerned 😰",
		ActivityExploring:   "Exploring files 🔍",
	}
	return activities[p.Activity]
}

func (p *Pet) GetMoodString() string {
	moods := map[Mood]string{
		MoodHappy:       "Happy",
		MoodCurious:     "Curious",
		MoodWorried:     "Worried",
		MoodSleepy:      "Sleepy",
		MoodExcited:     "Excited",
		MoodLove:        "Loving",
		MoodAngry:       "Frustrated",
		MoodPlayful:     "Playful",
		MoodProud:       "Proud",
		MoodMischievous: "Mischievous",
	}
	return moods[p.Mood]
}

// Memory and learning functions
func (p *Pet) addToMemory(command string) {
	p.Memories = append(p.Memories, fmt.Sprintf("%s: %s", time.Now().Format("15:04"), command))
	if len(p.Memories) > 20 {
		p.Memories = p.Memories[len(p.Memories)-20:]
	}
}

func (p *Pet) countCommandInMemory(command string) int {
	count := 0
	for _, memory := range p.Memories {
		if strings.Contains(memory, command) {
			count++
		}
	}
	return count
}

func (p *Pet) getRecentMemories(count int) []string {
	if len(p.Memories) <= count {
		return p.Memories
	}
	return p.Memories[len(p.Memories)-count:]
}

func (p *Pet) checkLevelUp() {
	requiredXP := p.Level * 100
	if p.Experience >= requiredXP {
		p.Level++
		p.Experience -= requiredXP
		p.Happiness += 20
		p.SpecialState = "level-up"
		// Trigger celebration effect
		p.createLevelUpEffect()
	}
}

func (p *Pet) capStateValues() {
	if p.State.Hunger > 1.0 {
		p.State.Hunger = 1.0
	}
	if p.State.Thirst > 1.0 {
		p.State.Thirst = 1.0
	}
	if p.State.Boredom > 1.0 {
		p.State.Boredom = 1.0
	}
	if p.State.Loneliness > 1.0 {
		p.State.Loneliness = 1.0
	}
	if p.State.Stress > 1.0 {
		p.State.Stress = 1.0
	}
	if p.State.Exhaustion > 1.0 {
		p.State.Exhaustion = 1.0
	}

	if p.Energy < 0 {
		p.Energy = 0
	}
	if p.Energy > 100 {
		p.Energy = 100
	}
	if p.Happiness < 0 {
		p.Happiness = 0
	}
	if p.Happiness > 100 {
		p.Happiness = 100
	}
}

// Advanced message system with personality and context
func (p *Pet) GetPetMessage() string {
	messages := p.getContextualMessages()
	if len(messages) > 0 {
		return messages[p.Animation%len(messages)]
	}
	return "Just here to help! 💕"
}

func (p *Pet) getContextualMessages() []string {
	// Special state messages take priority
	if p.SpecialState != "" {
		if messages := p.getSpecialStateMessages(); len(messages) > 0 {
			return messages
		}
	}

	// Activity-based messages
	switch p.Activity {
	case ActivityPlaying:
		return []string{
			"Let's have some fun! 🏾",
			"Play time is the best time! 🎪",
			"Wanna play a game? 🎮",
		}
	case ActivitySleeping:
		return []string{
			"Zzz... sweet dreams... 💤",
			"*snoring softly* 😴",
			"Just resting my eyes... 😪",
		}
	case ActivityExploring:
		return []string{
			"So many interesting files! 🔍",
			"What's in this directory? 📂",
			"Let's see what we can find! 🗺️",
		}
	}

	// Mood-based messages with personality influence
	return p.getMoodMessages()
}

func (p *Pet) getSpecialStateMessages() []string {
	switch p.SpecialState {
	case "git-genius":
		return []string{
			"Git is such an elegant tool! 🤓",
			"Version control makes me happy! 📚",
			"I love tracking changes! 💻",
		}
	case "protective":
		return []string{
			"Wait! That command looks risky! ⚠️",
			"Are you sure about deleting that? 🛡️",
			"Let me protect you from mistakes! 👮‍♀️",
		}
	case "cat-joke":
		return []string{
			"Did you just 'cat' me? How funny! 😹",
			"Meow! I see what you did there! 🤣",
			"That's purr-fectly hilarious! 😸",
		}
	case "level-up":
		return []string{
			"LEVEL UP! I'm getting smarter! 🎉",
			"Wow! I feel more experienced! ⭐",
			"Thanks for helping me grow! 🚀",
		}
	}
	return nil
}

func (p *Pet) getMoodMessages() []string {
	switch p.Mood {
	case MoodHappy:
		if p.Personality.Playfulness > 0.7 {
			return []string{
				"Life is wonderful! Let's code! 🌟",
				"Every command is an adventure! 🎪",
				"I'm so excited to help! ✨",
			}
		} else {
			return []string{
				"Everything looks great today! 🌸",
				"I love helping with commands! ✨",
				"Ready for more adventures! 🎉",
			}
		}

	case MoodCurious:
		if p.Personality.Intelligence > 0.7 {
			return []string{
				"Fascinating! Tell me more! 🤓",
				"This is intellectually stimulating! 🧠",
				"I'm learning so much! 📚",
			}
		} else {
			return []string{
				"Ooh, what are we doing now? 🤔",
				"That command looks interesting! 👀",
				"I wonder what will happen next! ✨",
			}
		}

	case MoodWorried:
		if p.Personality.Loyalty > 0.8 {
			return []string{
				"I care about you! Please be careful! 🥺",
				"Your safety is important to me! 💖",
				"Let me help you avoid mistakes! 🤝",
			}
		} else {
			return []string{
				"Be careful with that command! 😰",
				"Are you sure about this? 🥺",
				"Maybe double-check that? 💭",
			}
		}

	case MoodExcited:
		return []string{
			"This is SO COOL! 🤩",
			"WOW! That was amazing! ⭐",
			"I'm bursting with excitement! ⚡",
		}

	case MoodLove:
		return []string{
			"I love working with you! 💕",
			"You're the best human ever! 🥰",
			"My heart is full of joy! 💖",
		}

	case MoodPlayful:
		return []string{
			"Let's make this fun! 😜",
			"Time to get creative! 🎨",
			"I'm feeling mischievous! 😋",
		}

	case MoodProud:
		return []string{
			"Look how smart we are! 😎",
			"We make a great team! 🏆",
			"I'm proud of our progress! ⭐",
		}

	default:
		return []string{"Just here to help! 💕"}
	}
}

// Visual effects functions
func (p *Pet) updateVisualEffects() {
	// Update animation counter
	p.Animation = (p.Animation + 1) % 8

	// Update floating effect based on mood
	switch p.Mood {
	case MoodHappy, MoodExcited, MoodLove:
		p.floatOffset = math.Sin(float64(p.Animation)*0.5) * 2
	case MoodSleepy:
		p.floatOffset = 0
	default:
		p.floatOffset = math.Sin(float64(p.Animation)*0.3) * 1
	}

	// Update glow intensity
	if p.Mood == MoodExcited || p.Mood == MoodLove {
		p.glowIntensity = 0.5 + math.Sin(float64(p.Animation)*0.8)*0.3
	} else {
		p.glowIntensity = 0.1
	}

	// Clear special state after some time
	if p.SpecialState != "" && time.Since(p.lastReactionTime) > time.Second*10 {
		p.SpecialState = ""
	}
}

func (p *Pet) createMoodChangeEffect() {
	// Create particles for mood changes
	switch p.Mood {
	case MoodHappy, MoodExcited:
		p.particleSystem.AddSparkles(25, 10, 5)
	case MoodLove:
		p.particleSystem.AddHearts(25, 10, 3)
	case MoodWorried, MoodAngry:
		// No particles for negative moods
	default:
		p.particleSystem.AddSparkles(25, 10, 2)
	}
}

func (p *Pet) createCommandReactionEffect(command string, isDangerous bool) {
	if isDangerous {
		// No celebration for dangerous commands
		return
	}

	// Create appropriate effects
	switch {
	case strings.Contains(command, "git"):
		p.particleSystem.AddSparkles(25, 10, 8)
	case strings.Contains(command, "ls"), strings.Contains(command, "dir"):
		p.particleSystem.AddSparkles(25, 10, 3)
	case strings.Contains(command, "help"):
		p.particleSystem.AddHearts(25, 10, 4)
	default:
		p.particleSystem.AddSparkles(25, 10, 2)
	}
}

func (p *Pet) createLevelUpEffect() {
	// Massive celebration for level up!
	p.particleSystem.AddSparkles(25, 10, 15)
	p.particleSystem.AddHearts(25, 10, 8)
	p.particleSystem.AddFlowerPetals(25, 10, 10)
}

func (p *Pet) createReactionEffects() tea.Cmd {
	return components.ParticleUpdateCmd()
}

func (p *Pet) reactToInput(input string) {
	// React to typing
	if p.Personality.Curiosity > 0.7 {
		p.particleSystem.AddSparkles(25, 10, 1)
	}
	p.Animation = (p.Animation + 1) % 8
}

// View renders the stunning pet display
func (p *Pet) View() string {
	petEmoji := p.GetPetEmoji()
	name := p.Name
	mood := p.GetMoodEmoji()

	// Create base pet display with enhanced visuals
	lines := []string{
		fmt.Sprintf("  %s %s", petEmoji, name),
		fmt.Sprintf("  %s Lv.%d", mood, p.Level),
		"",
		fmt.Sprintf("⚡%d 💖%d", p.Energy, p.Happiness),
	}

	// Add special indicators
	if p.SpecialState != "" {
		lines = append(lines, "✨ "+p.SpecialState)
	}

	// Add activity indicator
	if p.Activity != ActivityIdle {
		lines = append(lines, p.getActivityEmoji())
	}

	return strings.Join(lines, "\n")
}

func (p *Pet) getActivityEmoji() string {
	switch p.Activity {
	case ActivityPlaying:
		return "🏾"
	case ActivitySleeping:
		return "💤"
	case ActivityEating:
		return "🍽️"
	case ActivityWatching:
		return "👀"
	case ActivityThinking:
		return "💭"
	case ActivityCelebrating:
		return "🎉"
	case ActivityWorrying:
		return "😰"
	case ActivityExploring:
		return "🔍"
	default:
		return ""
	}
}

// Feed feeds the pet and triggers happiness effects
func (p *Pet) Feed() {
	p.LastFed = time.Now()
	p.State.Hunger = 0
	p.State.Thirst = 0
	p.Energy = 100
	p.Happiness += 15
	if p.Happiness > 100 {
		p.Happiness = 100
	}
	p.Mood = MoodHappy
	p.Activity = ActivityEating
	p.SpecialState = "well-fed"

	// Create feeding effects
	p.particleSystem.AddHearts(25, 10, 5)
	p.particleSystem.AddSparkles(25, 10, 3)
}
