package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiiiiiis/hitori-tui/internal/api"
	hauth "github.com/mathiiiiiis/hitori-tui/internal/auth"
)

// interact menu actions
var interactActions = []string{"Feed", "Play", "Bathe", "Sleep / Wake", "Solve problem"}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case animTickMsg:
		m.frame = 1 - m.frame
		return m, animTick()

	case lifeTickMsg:
		if m.state != nil {
			m.state.Tick()
		}
		return m, lifeTick()

	case syncTickMsg:
		cmds := []tea.Cmd{syncTick()}
		if m.client != nil && m.state != nil {
			cmds = append(cmds, m.saveState())
		}
		return m, tea.Batch(cmds...)

	case authInitMsg:
		m.sessionID = msg.sessionID
		m.authURL = msg.authURL
		hauth.OpenBrowser(msg.authURL)
		return m, pollAuth(m.sessionID)

	case authPollMsg:
		if msg.err != nil {
			m.syncErr = msg.err.Error()
			return m, nil
		}
		if msg.token == "" {
			return m, pollAuth(m.sessionID)
		}
		hauth.StoreToken(msg.token)
		m.token = msg.token
		m.client = api.New(msg.token)
		m.view = viewWorld
		return m, tea.Batch(m.loadState(), lifeTick(), syncTick())

	case stateLoadedMsg:
		if msg.state == nil {
			// no Mono yet > creation wizard
			m.view = viewCreate
			m.create = newCreateModel()
			return m, textinput.Blink
		}
		m.state = msg.state
		m.synced = true
		m.view = viewWorld
		return m, nil

	case syncDoneMsg:
		if msg.err != nil {
			m.syncErr = msg.err.Error()
		} else {
			m.synced = true
			m.syncErr = ""
		}
		return m, nil

	case errMsg:
		m.syncErr = msg.err.Error()
		return m, nil
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// global quit
	if msg.String() == "ctrl+c" {
		if m.client != nil && m.state != nil {
			m.client.SaveState(m.state)
		}
		return m, tea.Quit
	}

	switch m.view {
	case viewCreate:
		return m.handleCreateKey(msg)
	case viewInteract:
		return m.handleInteractKey(msg)
	case viewCustomize:
		return m.handleCustomizeKey(msg)
	case viewWorld:
		return m.handleWorldKey(msg)
	}
	return m, nil
}

// ==== world ====

func (m *Model) handleWorldKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		if m.client != nil && m.state != nil {
			m.client.SaveState(m.state)
		}
		return m, tea.Quit
	case "i", "enter":
		m.view = viewInteract
		m.interactCursor = 0
	case "c":
		m.view = viewCustomize
		m.customCursor = 0
	case "s":
		return m, m.saveState()
	}
	return m, nil
}

// ==== interact ====

func (m *Model) handleInteractKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "i":
		m.view = viewWorld
	case "up", "k":
		if m.interactCursor > 0 {
			m.interactCursor--
		}
	case "down", "j":
		if m.interactCursor < len(interactActions)-1 {
			m.interactCursor++
		}
	case "enter", " ":
		m.doInteraction(m.interactCursor)
	}
	return m, nil
}

func (m *Model) doInteraction(idx int) {
	if m.state == nil {
		return
	}
	switch idx {
	case 0:
		m.state.Feed("")
	case 1:
		m.state.Play()
	case 2:
		m.state.Bathe()
	case 3:
		m.state.Sleep()
	case 4:
		m.state.SolveProblem()
	}
}

// ==== customize ====

func (m *Model) handleCustomizeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == nil {
		m.view = viewWorld
		return m, nil
	}
	switch msg.String() {
	case "esc", "q", "c":
		m.view = viewWorld
		return m, m.saveState()
	case "up", "k":
		if m.customCursor > 0 {
			m.customCursor--
		}
	case "down", "j":
		if m.customCursor < 4 {
			m.customCursor++
		}
	case "left", "h":
		m.adjustPersonality(m.customCursor, -5)
	case "right", "l":
		m.adjustPersonality(m.customCursor, +5)
	}
	return m, nil
}

func (m *Model) adjustPersonality(idx, delta int) {
	p := &m.state.Personality
	vals := []*float64{&p.Patience, &p.EnergyLevel, &p.Pickiness, &p.Sociability, &p.Optimism}
	if idx < 0 || idx >= len(vals) {
		return
	}
	v := *vals[idx] + float64(delta)
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	*vals[idx] = v
}

// ==== creation wizard ====

func (m *Model) handleCreateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	c := &m.create

	switch msg.String() {
	case "esc":
		// step back
		if c.step > stepName {
			c.step--
			c.focusStep()
		}
		return m, nil

	case "enter":
		if c.step == stepConfirm {
			// finalize
			m.state = c.build()
			m.state.SimulateOffline()
			m.view = viewWorld
			return m, m.saveState()
		}
		if c.step == stepName && c.name.Value() == "" {
			return m, nil // name required
		}
		c.step++
		c.focusStep()
		return m, nil
	}

	// personality step: arrow keys adjust sliders
	if c.step == stepPersonality {
		switch msg.String() {
		case "up", "k":
			if c.pIndex > 0 {
				c.pIndex--
			}
		case "down", "j":
			if c.pIndex < 4 {
				c.pIndex++
			}
		case "left", "h":
			c.personality[c.pIndex] = clampInt(c.personality[c.pIndex]-5, 0, 100)
		case "right", "l":
			c.personality[c.pIndex] = clampInt(c.personality[c.pIndex]+5, 0, 100)
		}
		return m, nil
	}

	// text steps: forward to active input
	var cmd tea.Cmd
	switch c.step {
	case stepName:
		c.name, cmd = c.name.Update(msg)
	case stepBirthday:
		c.birthday, cmd = c.birthday.Update(msg)
	case stepCatchphrase:
		c.catchphrase, cmd = c.catchphrase.Update(msg)
	}
	return m, cmd
}

