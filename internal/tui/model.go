package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiiiiiis/hitori-tui/internal/api"
	"github.com/mathiiiiiis/hitori-tui/internal/mono"
)

// ==== view states ====

type viewState int

const (
	viewLogin viewState = iota
	viewCreate
	viewWorld
	viewInteract
	viewCustomize
)

// ==== tea messages ====

type animTickMsg struct{}
type lifeTickMsg struct{}
type syncTickMsg struct{}
type authInitMsg struct {
	sessionID string
	authURL   string
}
type authPollMsg struct {
	token string
	err   error
}
type stateLoadedMsg struct{ state *mono.State } // nil => no Mono yet, go to creation
type syncDoneMsg struct{ err error }
type errMsg struct{ err error }

// ==== model ====

type Model struct {
	width  int
	height int
	view   viewState

	// auth
	sessionID string
	authURL   string
	spinner   spinner.Model

	// one Mono
	state *mono.State
	frame int

	// backend
	client  *api.Client
	token   string
	synced  bool
	syncErr string

	// creation wizard
	create createModel

	// interact menu
	interactCursor int

	// customize editor
	customCursor int
}

func New(token string) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle

	m := &Model{
		spinner: sp,
		token:   token,
		create:  newCreateModel(),
	}

	if token != "" {
		m.client = api.New(token)
		m.view = viewWorld // will redirect to create if no Mono
	} else {
		m.view = viewLogin
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick, animTick()}
	if m.view == viewLogin {
		cmds = append(cmds, m.initCLIAuth())
	} else {
		cmds = append(cmds, m.loadState(), lifeTick(), syncTick())
	}
	return tea.Batch(cmds...)
}

// ==== ticks ====

func animTick() tea.Cmd {
	return tea.Tick(700*time.Millisecond, func(time.Time) tea.Msg { return animTickMsg{} })
}

func lifeTick() tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg { return lifeTickMsg{} })
}

func syncTick() tea.Cmd {
	return tea.Tick(60*time.Second, func(time.Time) tea.Msg { return syncTickMsg{} })
}

// ==== async commands ====

func (m *Model) initCLIAuth() tea.Cmd {
	client := api.New("")
	return func() tea.Msg {
		init, err := client.CLIAuthInit("discord")
		if err != nil {
			return errMsg{err}
		}
		return authInitMsg{init.SessionID, init.AuthURL}
	}
}

func pollAuth(sessionID string) tea.Cmd {
	client := api.New("")
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		token, err := client.CLIAuthPoll(sessionID)
		return authPollMsg{token: token, err: err}
	})
}

func (m *Model) loadState() tea.Cmd {
	return func() tea.Msg {
		s, err := m.client.LoadState()
		if err != nil {
			return errMsg{err}
		}
		if s == nil {
			return stateLoadedMsg{state: nil} // no Mono => creation
		}
		if s.Level == 0 {
			s.Level = 1
		}
		s.SimulateOffline()
		return stateLoadedMsg{state: s}
	}
}

func (m *Model) saveState() tea.Cmd {
	state := m.state
	client := m.client
	return func() tea.Msg {
		if client == nil || state == nil {
			return syncDoneMsg{}
		}
		state.SavedAt = time.Now().UnixMilli()
		return syncDoneMsg{client.SaveState(state)}
	}
}

// newTextInput is a small helper used by creation wizard
func newTextInput(placeholder string, limit int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = limit
	ti.Prompt = "> "
	ti.PromptStyle = promptStyle
	return ti
}
