package mono

import "github.com/charmbracelet/lipgloss"

var sprites = map[Mood][2][5]string{
	MoodHappy: {
		{
			"   .---.   ",
			"  ( ^ ^ )  ",
			"  |  ◡  |  ",
			"   '---'   ",
			"   /| |\\   ",
		},
		{
			"   .---.   ",
			"  ( ^ ^ )  ",
			"  |  ◡  |  ",
			"   '---'   ",
			"   \\| |/   ",
		},
	},
	MoodContent: {
		{
			"   .---.   ",
			"  ( • • )  ",
			"  |  ‿  |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( - • )  ",
			"  |  ‿  |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodNeutral: {
		{
			"   .---.   ",
			"  ( - - )  ",
			"  |  -  |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( • • )  ",
			"  |  -  |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodSad: {
		{
			"   .---.   ",
			"  ( ╥ ╥ )  ",
			"  |  ︵ |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( ╥ ╥ )  ",
			"  |  ︵ |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodBored: {
		{
			"   .---.   ",
			"  ( ─ ─ )  ",
			"  |  ~  |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( ─ • )  ",
			"  |  ~  | .",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodTired: {
		{
			"   .---.   ",
			"  ( ˘ ˘ )  ",
			"  |  ‿  |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( ¬ ¬ )  ",
			"  |  ‿  |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodExhausted: {
		{
			"   .---.   ",
			"  ( × × )  ",
			"  |  ﹏ |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( × × )  ",
			"  |  ﹏ |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodStarving: {
		{
			"   .---.   ",
			"  ( ｡ ｡ )  ",
			"  |  ○  |  ",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"   ( ｡ ｡ )  ",
			"  |   o |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
	MoodSick: {
		{
			"   .---.   ",
			"  ( @ @ )  ",
			"  |  ~  | +",
			"   '---'   ",
			"    | |    ",
		},
		{
			"   .---.   ",
			"  ( @ @ )  ",
			"  |  ~  |  ",
			"   '---'   ",
			"    | |    ",
		},
	},
}

var sleeping = [2][5]string{
	{
		"   .---.  z",
		"  ( - - )  ",
		"  |  ‿  |  ",
		"   '---'   ",
		"    | |    ",
	},
	{
		"   .---. zZ",
		"  ( _ _ )  ",
		"  |  ‿  |  ",
		"   '---'   ",
		"    | |    ",
	},
}

// moodColor maps each mood to a body color
var moodColor = map[Mood]lipgloss.Color{
	MoodHappy:     lipgloss.Color("#AADDFF"),
	MoodContent:   lipgloss.Color("#688BFF"),
	MoodNeutral:   lipgloss.Color("#7A7AAA"),
	MoodSad:       lipgloss.Color("#5B7BC4"),
	MoodBored:     lipgloss.Color("#6B6B99"),
	MoodTired:     lipgloss.Color("#7766AA"),
	MoodExhausted: lipgloss.Color("#9955AA"),
	MoodStarving:  lipgloss.Color("#FFAA44"),
	MoodSick:      lipgloss.Color("#88CC66"),
}

// RenderSprite returns the colored sprite for a mood + frame
func RenderSprite(mood Mood, frame int, asleep bool) string {
	if frame < 0 || frame > 1 {
		frame = 0
	}

	var lines [5]string
	var color lipgloss.Color

	if asleep {
		lines = sleeping[frame]
		color = lipgloss.Color("#445588")
	} else {
		f, ok := sprites[mood]
		if !ok {
			f = sprites[MoodNeutral]
		}
		lines = f[frame]
		color = moodColor[mood]
		if color == "" {
			color = lipgloss.Color("#688BFF")
		}
	}

	style := lipgloss.NewStyle().Foreground(color)
	out := ""
	for i, line := range lines {
		out += style.Render(line)
		if i < 4 {
			out += "\n"
		}
	}
	return out
}
