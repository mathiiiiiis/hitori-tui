package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mathiiiiiis/hitori-tui/internal/mono"
)

type createModel struct {
	step int // which field is active

	name        textinput.Model
	birthday    textinput.Model
	catchphrase textinput.Model

	// personality sliders, 0–100
	personality [5]int
	pIndex      int // which slider is selected on personality step
}

// creation steps
const (
	stepName = iota
	stepBirthday
	stepCatchphrase
	stepPersonality
	stepConfirm
	stepCount
)

var personalityLabels = [5]string{"Patience", "Energy", "Pickiness", "Sociability", "Optimism"}

func newCreateModel() createModel {
	c := createModel{
		name:        newTextInput("Mono's name", 20),
		birthday:    newTextInput("MM-DD (optional)", 5),
		catchphrase: newTextInput("a catchphrase (optional)", 40),
		personality: [5]int{50, 50, 50, 50, 50},
	}
	c.name.Focus()
	return c
}

// build assembles fresh State from wizard inputs
func (c createModel) build() *mono.State {
	s := mono.NewState()
	s.Name = c.name.Value()
	s.Birthday = c.birthday.Value()
	s.Catchphrase = c.catchphrase.Value()
	s.Personality = mono.Personality{
		Patience:    float64(c.personality[0]),
		EnergyLevel: float64(c.personality[1]),
		Pickiness:   float64(c.personality[2]),
		Sociability: float64(c.personality[3]),
		Optimism:    float64(c.personality[4]),
	}
	return s
}

// focusStep moves focus to active text input for current step
func (c *createModel) focusStep() {
	c.name.Blur()
	c.birthday.Blur()
	c.catchphrase.Blur()
	switch c.step {
	case stepName:
		c.name.Focus()
	case stepBirthday:
		c.birthday.Focus()
	case stepCatchphrase:
		c.catchphrase.Focus()
	}
}
