package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ComponentType represents different types of UI components
type ComponentType int

const (
	ComponentButton ComponentType = iota
	ComponentSlider
	ComponentProgressBar
	ComponentTabGroup
	ComponentDropdown
	ComponentModal
)

// Component represents a generic UI component interface
type Component interface {
	Update(msg tea.Msg) (Component, tea.Cmd)
	Render() string
	Focus()
	Blur()
}

// ComponentManager manages all interactive UI components
type ComponentManager struct {
	components    []Component
	focusedIndex  int
	globalEffects *ParticleSystem
	animation     *AnimationManager
	lastUpdate    time.Time
}

// NewComponentManager creates a new component manager
func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components:    make([]Component, 0),
		focusedIndex:  -1,
		globalEffects: NewParticleSystem(100, 50),
		animation:     NewAnimationManager(),
		lastUpdate:    time.Now(),
	}
}

// AddComponent adds a component to the manager
func (cm *ComponentManager) AddComponent(component Component) {
	cm.components = append(cm.components, component)

	// Focus first component if none focused
	if cm.focusedIndex == -1 {
		cm.focusedIndex = 0
		component.Focus()
	}
}

// Update updates all components
func (cm *ComponentManager) Update(msg tea.Msg) (*ComponentManager, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			cm.NextComponent()
		case "shift+tab":
			cm.PrevComponent()
		}

	case ParticleTickMsg:
		cm.globalEffects.Update(0.05)
		cm.animation.Update()
		cmds = append(cmds, ParticleUpdateCmd())
	}

	// Update all components
	for i, component := range cm.components {
		updatedComponent, cmd := component.Update(msg)
		cm.components[i] = updatedComponent
		cmds = append(cmds, cmd)
	}

	return cm, tea.Batch(cmds...)
}

// NextComponent focuses the next component
func (cm *ComponentManager) NextComponent() {
	if len(cm.components) == 0 {
		return
	}

	// Blur current component
	if cm.focusedIndex >= 0 && cm.focusedIndex < len(cm.components) {
		cm.components[cm.focusedIndex].Blur()
	}

	// Move to next component
	cm.focusedIndex = (cm.focusedIndex + 1) % len(cm.components)

	// Focus new component
	cm.components[cm.focusedIndex].Focus()

	// Create focus change effect
	cm.globalEffects.AddSparkles(50, 25, 8)
}

// PrevComponent focuses the previous component
func (cm *ComponentManager) PrevComponent() {
	if len(cm.components) == 0 {
		return
	}

	// Blur current component
	if cm.focusedIndex >= 0 && cm.focusedIndex < len(cm.components) {
		cm.components[cm.focusedIndex].Blur()
	}

	// Move to previous component
	cm.focusedIndex = (cm.focusedIndex - 1 + len(cm.components)) % len(cm.components)

	// Focus new component
	cm.components[cm.focusedIndex].Focus()

	// Create focus change effect
	cm.globalEffects.AddSparkles(50, 25, 8)
}

// SetFocusedComponent sets focus to specific component
func (cm *ComponentManager) SetFocusedComponent(index int) {
	if index < 0 || index >= len(cm.components) {
		return
	}

	// Blur current component
	if cm.focusedIndex >= 0 && cm.focusedIndex < len(cm.components) {
		cm.components[cm.focusedIndex].Blur()
	}

	// Set new focused component
	cm.focusedIndex = index
	cm.components[index].Focus()

	// Create focus change effect
	cm.globalEffects.AddSparkles(50, 25, 8)
}

// GetFocusedComponent returns the currently focused component
func (cm *ComponentManager) GetFocusedComponent() Component {
	if cm.focusedIndex >= 0 && cm.focusedIndex < len(cm.components) {
		return cm.components[cm.focusedIndex]
	}
	return nil
}

// RemoveComponent removes a component from the manager
func (cm *ComponentManager) RemoveComponent(index int) {
	if index < 0 || index >= len(cm.components) {
		return
	}

	// Blur the component being removed
	cm.components[index].Blur()

	// Remove component
	cm.components = append(cm.components[:index], cm.components[index+1:]...)

	// Adjust focused index
	if cm.focusedIndex == index {
		if len(cm.components) > 0 {
			if cm.focusedIndex >= len(cm.components) {
				cm.focusedIndex = len(cm.components) - 1
			}
			cm.components[cm.focusedIndex].Focus()
		} else {
			cm.focusedIndex = -1
		}
	} else if cm.focusedIndex > index {
		cm.focusedIndex--
	}
}

// CreateGlobalEffect creates a global particle effect
func (cm *ComponentManager) CreateGlobalEffect(effectType string, x, y, count int) {
	switch effectType {
	case "sparkles":
		cm.globalEffects.AddSparkles(x, y, count)
	case "hearts":
		cm.globalEffects.AddHearts(x, y, count)
	case "flowers":
		cm.globalEffects.AddFlowerPetals(x, y, count)
	}
}

// Render renders all components (this would be used by parent UI)
func (cm *ComponentManager) Render() string {
	// This method would be implemented based on layout requirements
	// For now, just return empty string as components are rendered individually
	return ""
}

// ButtonWrapper wraps Button to implement Component interface
type ButtonWrapper struct {
	*Button
}

func (bw *ButtonWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	button, cmd := bw.Button.Update(msg)
	bw.Button = button
	return bw, cmd
}

func (bw *ButtonWrapper) Focus() {
	bw.Button.Focus()
}

func (bw *ButtonWrapper) Blur() {
	bw.Button.Blur()
}

// SliderWrapper wraps Slider to implement Component interface
type SliderWrapper struct {
	*Slider
}

func (sw *SliderWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	slider, cmd := sw.Slider.Update(msg)
	sw.Slider = slider
	return sw, cmd
}

func (sw *SliderWrapper) Focus() {
	sw.Slider.Focus()
}

func (sw *SliderWrapper) Blur() {
	sw.Slider.Blur()
}

// ProgressBarWrapper wraps ProgressBar to implement Component interface
type ProgressBarWrapper struct {
	*ProgressBar
}

func (pbw *ProgressBarWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	pb, cmd := pbw.ProgressBar.Update(msg)
	pbw.ProgressBar = pb
	return pbw, cmd
}

func (pbw *ProgressBarWrapper) Focus() {
	// ProgressBar doesn't need focus functionality but we implement it
}

func (pbw *ProgressBarWrapper) Blur() {
	// ProgressBar doesn't need blur functionality but we implement it
}

// TabGroupWrapper wraps TabGroup to implement Component interface
type TabGroupWrapper struct {
	*TabGroup
}

func (tgw *TabGroupWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	tabGroup, cmd := tgw.TabGroup.Update(msg)
	tgw.TabGroup = tabGroup
	return tgw, cmd
}

func (tgw *TabGroupWrapper) Focus() {
	tgw.TabGroup.Focus()
}

func (tgw *TabGroupWrapper) Blur() {
	tgw.TabGroup.Blur()
}

// DropdownWrapper wraps Dropdown to implement Component interface
type DropdownWrapper struct {
	*Dropdown
}

func (dw *DropdownWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	dropdown, cmd := dw.Dropdown.Update(msg)
	dw.Dropdown = dropdown
	return dw, cmd
}

func (dw *DropdownWrapper) Focus() {
	dw.Dropdown.Focus()
}

func (dw *DropdownWrapper) Blur() {
	dw.Dropdown.Blur()
}

// ModalWrapper wraps Modal to implement Component interface
type ModalWrapper struct {
	*Modal
}

func (mw *ModalWrapper) Update(msg tea.Msg) (Component, tea.Cmd) {
	modal, cmd := mw.Modal.Update(msg)
	mw.Modal = modal
	return mw, cmd
}

func (mw *ModalWrapper) Focus() {
	mw.Modal.Focus()
}

func (mw *ModalWrapper) Blur() {
	mw.Modal.Blur()
}

// Helper functions to create wrapped components

// WrapButton wraps a Button to implement Component interface
func WrapButton(button *Button) Component {
	return &ButtonWrapper{Button: button}
}

// WrapSlider wraps a Slider to implement Component interface
func WrapSlider(slider *Slider) Component {
	return &SliderWrapper{Slider: slider}
}

// WrapProgressBar wraps a ProgressBar to implement Component interface
func WrapProgressBar(pb *ProgressBar) Component {
	return &ProgressBarWrapper{ProgressBar: pb}
}

// WrapTabGroup wraps a TabGroup to implement Component interface
func WrapTabGroup(tabGroup *TabGroup) Component {
	return &TabGroupWrapper{TabGroup: tabGroup}
}

// WrapDropdown wraps a Dropdown to implement Component interface
func WrapDropdown(dropdown *Dropdown) Component {
	return &DropdownWrapper{Dropdown: dropdown}
}

// WrapModal wraps a Modal to implement Component interface
func WrapModal(modal *Modal) Component {
	return &ModalWrapper{Modal: modal}
}
