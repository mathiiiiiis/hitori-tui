package mono

import "time"

type Personality struct {
	Patience    float64 `json:"patience"`
	EnergyLevel float64 `json:"energyLevel"`
	Pickiness   float64 `json:"pickiness"`
	Sociability float64 `json:"sociability"`
	Optimism    float64 `json:"optimism"`
}

type Needs struct {
	Hunger    float64 `json:"hunger"`
	Energy    float64 `json:"energy"`
	Happiness float64 `json:"happiness"`
	Hygiene   float64 `json:"hygiene"`
	Fun       float64 `json:"fun"`
	Health    float64 `json:"health"`
}

type State struct {
	// ==== identity ====
	Name        string `json:"name"`
	Birthday    string `json:"birthday"`
	Catchphrase string `json:"catchphrase"`

	// ==== customization ====
	Personality Personality `json:"personality"`

	// ==== live state ====
	Needs    Needs  `json:"needs"`
	Activity string `json:"activity"` // idle | sleeping | napping | eating | playing

	// ==== progress ====
	Level          int           `json:"level"`
	XP             int           `json:"xp"`
	Inventory      []string      `json:"inventory"`
	Room           map[string]any `json:"room"`
	Clothing       any           `json:"clothing"`
	FavoriteFood   any           `json:"favoriteFood"`
	WorstFood      any           `json:"worstFood"`
	ProblemsSolved int           `json:"problemsSolved"`

	// ==== events ====
	CurrentProblem any      `json:"currentProblem"`
	DreamFlags     []string `json:"dreamFlags"`
	EventHistory   []string `json:"eventHistory"`

	// ==== meta ====
	SavedAt int64 `json:"savedAt"` // unix ms; used by both sides for offline simulation
}

// NewState returns an uninitialized Mono
func NewState() *State {
	return &State{
		Personality: Personality{50, 50, 50, 50, 50},
		Needs:       Needs{100, 100, 100, 100, 100, 100},
		Activity:    "idle",
		Level:       1,
		XP:          0,
		Inventory:   []string{},
		Room:        map[string]any{},
		DreamFlags:  []string{},
		EventHistory: []string{},
		SavedAt:     time.Now().UnixMilli(),
	}
}

// Mood is derived emotional state
type Mood string

const (
	MoodSick      Mood = "sick"
	MoodExhausted Mood = "exhausted"
	MoodTired     Mood = "tired"
	MoodStarving  Mood = "starving"
	MoodBored     Mood = "bored"
	MoodSad       Mood = "sad"
	MoodHappy     Mood = "happy"
	MoodContent   Mood = "content"
	MoodNeutral   Mood = "neutral"
)

// Mood derives current mood from needs + personality
func (s *State) Mood() Mood {
	n := s.Needs
	p := s.Personality

	switch {
	case n.Health < 20:
		return MoodSick
	case n.Energy < 10:
		return MoodExhausted
	case n.Energy < 30:
		return MoodTired
	case n.Hunger < 20:
		return MoodStarving
	case n.Fun < 20 && n.Happiness < 30:
		return MoodBored
	case n.Happiness < 20:
		return MoodSad
	}

	avg := (n.Hunger + n.Energy + n.Happiness + n.Hygiene + n.Fun + n.Health) / 6
	boosted := avg + (p.Optimism-50)*0.2

	switch {
	case boosted > 80:
		return MoodHappy
	case boosted > 50:
		return MoodContent
	default:
		return MoodNeutral
	}
}

func (s *State) IsAsleep() bool {
	return s.Activity == "sleeping" || s.Activity == "napping"
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
