package tui

import "github.com/charmbracelet/lipgloss"

const (
	ColorBlue    = lipgloss.Color("#688BFF")
	ColorPink    = lipgloss.Color("#FF4893")
	ColorSurface = lipgloss.Color("#1A1A28")
	ColorBorder  = lipgloss.Color("#2E2E44")
	ColorMuted   = lipgloss.Color("#6666AA")
	ColorText    = lipgloss.Color("#E0E0F0")
	ColorSubtle  = lipgloss.Color("#9090BB")
	ColorGreen   = lipgloss.Color("#44DD88")
	ColorOrange  = lipgloss.Color("#FF8844")
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(ColorBlue)
	promptStyle  = lipgloss.NewStyle().Foreground(ColorBlue)

	titleBarStyle = lipgloss.NewStyle().Foreground(ColorText).Background(ColorSurface).Padding(0, 2)
	footerStyle   = lipgloss.NewStyle().Foreground(ColorMuted).Background(ColorSurface).Padding(0, 2)

	monoPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(ColorBlue).Padding(1, 2)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(ColorBorder).Padding(1, 2)

	nameStyle   = lipgloss.NewStyle().Foreground(ColorBlue).Bold(true)
	levelStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	speechStyle = lipgloss.NewStyle().Foreground(ColorSubtle).Italic(true)
	moodStyle   = lipgloss.NewStyle().Foreground(ColorPink)

	barFilled  = lipgloss.NewStyle().Foreground(ColorBlue)
	barEmpty   = lipgloss.NewStyle().Foreground(ColorBorder)
	barLabel   = lipgloss.NewStyle().Foreground(ColorSubtle).Width(14)
	barPercent = lipgloss.NewStyle().Foreground(ColorMuted).Width(5)

	logStyle = lipgloss.NewStyle().Foreground(ColorSubtle)

	loginBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(ColorBlue).Padding(2, 4).Align(lipgloss.Center)
	createBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(ColorBlue).Padding(2, 4)

	urlStyle = lipgloss.NewStyle().Foreground(ColorBlue).Underline(true)

	keyStyle  = lipgloss.NewStyle().Foreground(ColorBlue).Bold(true)
	hintStyle = lipgloss.NewStyle().Foreground(ColorMuted)

	selectedStyle = lipgloss.NewStyle().Foreground(ColorPink).Bold(true)
	unselStyle    = lipgloss.NewStyle().Foreground(ColorSubtle)

	syncedStyle = lipgloss.NewStyle().Foreground(ColorGreen)
	errorStyle  = lipgloss.NewStyle().Foreground(ColorOrange)
)

// progressBar renders filled/empty bar for a 0-100 value
func progressBar(value, width int) string {
	filled := value * width / 100
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return barFilled.Render(repeat("█", filled)) + barEmpty.Render(repeat("░", width-filled))
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
