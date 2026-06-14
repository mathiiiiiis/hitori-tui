package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mathiiiiiis/hitori-tui/internal/mono"
)

func (m *Model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	switch m.view {
	case viewLogin:
		return m.loginView()
	case viewCreate:
		return m.createView()
	case viewInteract:
		return m.worldView(true)
	case viewCustomize:
		return m.customizeView()
	default:
		return m.worldView(false)
	}
}

// ==== login ====

func (m *Model) loginView() string {
	var b strings.Builder
	b.WriteString(nameStyle.Render("  hitori") + "\n")
	b.WriteString(levelStyle.Render("  a little life, in your terminal") + "\n\n")

	if m.authURL == "" {
		b.WriteString(m.spinner.View() + " connecting...\n")
	} else {
		b.WriteString("opening browser for Discord login...\n\n")
		b.WriteString("if it didn't open, visit:\n")
		b.WriteString(urlStyle.Render(m.authURL) + "\n\n")
		b.WriteString(m.spinner.View() + " waiting for authentication...\n")
	}
	if m.syncErr != "" {
		b.WriteString("\n" + errorStyle.Render("⚠ "+m.syncErr) + "\n")
	}
	b.WriteString("\n" + hintStyle.Render("ctrl+c to cancel"))

	box := loginBoxStyle.Width(clampInt(m.width-4, 40, 64)).Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// ==== creation wizard ====

func (m *Model) createView() string {
	c := m.create
	var b strings.Builder

	b.WriteString(nameStyle.Render("Create your Mono") + "\n")
	b.WriteString(levelStyle.Render(fmt.Sprintf("step %d of %d", c.step+1, stepCount)) + "\n\n")

	switch c.step {
	case stepName:
		b.WriteString("What's your Mono's name?\n\n")
		b.WriteString(c.name.View())
	case stepBirthday:
		b.WriteString("When's their birthday?\n\n")
		b.WriteString(c.birthday.View())
	case stepCatchphrase:
		b.WriteString("Give them a catchphrase.\n\n")
		b.WriteString(c.catchphrase.View())
	case stepPersonality:
		b.WriteString("Shape their personality.\n\n")
		for i, label := range personalityLabels {
			cursor := "  "
			style := unselStyle
			if i == c.pIndex {
				cursor = selectedStyle.Render("▸ ")
				style = selectedStyle
			}
			b.WriteString(cursor + style.Render(fmt.Sprintf("%-12s", label)) + " " +
				progressBar(c.personality[i], 14) + fmt.Sprintf(" %3d\n", c.personality[i]))
		}
		b.WriteString("\n" + hintStyle.Render("←/→ adjust · ↑/↓ select"))
	case stepConfirm:
		b.WriteString("Ready to bring them to life?\n\n")
		b.WriteString(nameStyle.Render(c.name.Value()) + "\n")
		if c.birthday.Value() != "" {
			b.WriteString(levelStyle.Render("birthday: "+c.birthday.Value()) + "\n")
		}
		if c.catchphrase.Value() != "" {
			b.WriteString(speechStyle.Render("\""+c.catchphrase.Value()+"\"") + "\n")
		}
		b.WriteString("\n")
		for i, label := range personalityLabels {
			b.WriteString(barLabel.Render(label) + " " + progressBar(c.personality[i], 12) +
				fmt.Sprintf(" %3d\n", c.personality[i]))
		}
	}

	b.WriteString("\n\n")
	if c.step == stepConfirm {
		b.WriteString(keyHint("enter", "create") + "   " + keyHint("esc", "back"))
	} else {
		b.WriteString(keyHint("enter", "next") + "   " + keyHint("esc", "back"))
	}

	box := createBoxStyle.Width(clampInt(m.width-4, 44, 64)).Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// ==== world ====

func (m *Model) worldView(interact bool) string {
	if m.state == nil {
		return "loading Mono..."
	}
	s := m.state

	// title bar
	title := nameStyle.Render("hitori") + levelStyle.Render(" · "+s.Name)
	badge := m.syncBadge()
	gap := clampInt(m.width-lipgloss.Width(title)-lipgloss.Width(badge)-4, 0, 200)
	titleBar := titleBarStyle.Width(m.width).Render(title + strings.Repeat(" ", gap) + badge)

	// footer
	var footer string
	if interact {
		footer = footerStyle.Width(m.width).Render(
			keyHint("↑/↓", "select") + "  " + keyHint("enter", "do") + "  " + keyHint("esc", "back"))
	} else {
		footer = footerStyle.Width(m.width).Render(
			keyHint("i", "interact") + "  " + keyHint("c", "customize") + "  " +
				keyHint("s", "sync") + "  " + keyHint("q", "quit"))
	}

	bodyH := m.height - 2
	leftW := 24
	rightW := m.width - leftW - 4

	left := m.monoPanel(leftW, bodyH)
	var right string
	if interact {
		right = m.interactPanel(rightW, bodyH)
	} else {
		right = m.statsPanel(rightW, bodyH)
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return titleBar + "\n" + body + "\n" + footer
}

func (m *Model) monoPanel(w, h int) string {
	s := m.state
	var b strings.Builder

	b.WriteString(mono.RenderSprite(s.Mood(), m.frame, s.IsAsleep()))
	b.WriteString("\n\n")

	if s.Catchphrase != "" && !s.IsAsleep() {
		b.WriteString(speechStyle.Render("\""+s.Catchphrase+"\"") + "\n")
	}
	b.WriteString(moodStyle.Render(string(s.Mood())))
	if s.Activity != "idle" {
		b.WriteString(levelStyle.Render("  · " + s.Activity))
	}

	return monoPanelStyle.Width(w).Height(h - 2).Render(b.String())
}

func (m *Model) statsPanel(w, h int) string {
	s := m.state
	n := s.Needs
	var b strings.Builder

	b.WriteString(nameStyle.Render(s.Name) + "  " + levelStyle.Render(fmt.Sprintf("Lv.%d", s.Level)) + "\n\n")

	b.WriteString(needBar("❤️  Happiness", n.Happiness) + "\n")
	b.WriteString(needBar("⚡  Energy   ", n.Energy) + "\n")
	b.WriteString(needBar("🍴  Hunger   ", n.Hunger) + "\n")
	b.WriteString(needBar("🚿  Hygiene  ", n.Hygiene) + "\n")
	b.WriteString(needBar("⭐  Fun      ", n.Fun) + "\n")
	b.WriteString(needBar("⛑️  Health   ", n.Health) + "\n\n")

	xpPct := 0
	if t := s.XPThreshold(); t > 0 {
		xpPct = s.XP * 100 / t
	}
	b.WriteString(levelStyle.Render("XP  ") + fmt.Sprintf("%d / %d  ", s.XP, s.XPThreshold()) +
		progressBar(xpPct, 8) + "\n\n")

	b.WriteString(levelStyle.Render("recent") + "\n")
	logs := s.EventHistory
	maxLines := clampInt(h-18, 2, 8)
	if len(logs) > maxLines {
		logs = logs[len(logs)-maxLines:]
	}
	for i := len(logs) - 1; i >= 0; i-- {
		b.WriteString(logStyle.Render(logs[i]) + "\n")
	}

	return panelStyle.Width(w).Height(h - 2).Render(b.String())
}

func (m *Model) interactPanel(w, h int) string {
	var b strings.Builder
	b.WriteString(nameStyle.Render("Interact") + "\n\n")

	for i, action := range interactActions {
		if i == m.interactCursor {
			b.WriteString(selectedStyle.Render("▸ "+action) + "\n")
		} else {
			b.WriteString(unselStyle.Render("  "+action) + "\n")
		}
	}

	b.WriteString("\n" + levelStyle.Render("Mono is " + string(m.state.Mood())))
	if m.state.IsAsleep() {
		b.WriteString("\n" + hintStyle.Render("(can't do much while asleep)"))
	}

	return panelStyle.Width(w).Height(h - 2).Render(b.String())
}

// ==== customize ====

func (m *Model) customizeView() string {
	s := m.state
	if s == nil {
		return "loading..."
	}
	p := s.Personality
	vals := []float64{p.Patience, p.EnergyLevel, p.Pickiness, p.Sociability, p.Optimism}

	titleBar := titleBarStyle.Width(m.width).Render(nameStyle.Render("hitori") + levelStyle.Render(" · Customize"))
	footer := footerStyle.Width(m.width).Render(
		keyHint("↑/↓", "select") + "  " + keyHint("←/→", "adjust") + "  " + keyHint("esc", "done"))

	var b strings.Builder
	b.WriteString(nameStyle.Render(s.Name) + "\n")
	if s.Catchphrase != "" {
		b.WriteString(speechStyle.Render("\""+s.Catchphrase+"\"") + "\n")
	}
	b.WriteString("\n" + levelStyle.Render("personality") + "\n\n")

	for i, label := range personalityLabels {
		cursor := "  "
		style := unselStyle
		if i == m.customCursor {
			cursor = selectedStyle.Render("▸ ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(fmt.Sprintf("%-12s", label)) + " " +
			progressBar(int(vals[i]), 16) + fmt.Sprintf(" %3.0f\n", vals[i]))
	}

	panel := panelStyle.Width(m.width - 4).Height(m.height - 4).Render(b.String())
	return titleBar + "\n" + panel + "\n" + footer
}

// ==== helpers ====

func (m *Model) syncBadge() string {
	if m.syncErr != "" {
		return errorStyle.Render("⚠ sync")
	}
	if m.synced {
		return syncedStyle.Render("● synced")
	}
	return levelStyle.Render("○ local")
}

func needBar(label string, value float64) string {
	v := int(value)
	return barLabel.Render(label) + " " + progressBar(v, 10) + " " +
		barPercent.Render(fmt.Sprintf("%3d%%", v))
}

func keyHint(key, desc string) string {
	return keyStyle.Render(key) + " " + hintStyle.Render(desc)
}
