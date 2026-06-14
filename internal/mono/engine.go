package mono

import "time"

// ==== tick / decay ====

// Tick advances life by one step
func (s *State) Tick() {
	n := &s.Needs
	p := s.Personality
	hour := time.Now().Hour()
	isNight := hour >= 22 || hour < 7

	// sleeping logic
	if s.Activity == "sleeping" || s.Activity == "napping" {
		n.Energy = clampF(n.Energy+3.5, 0, 100)
		n.Hunger = clampF(n.Hunger-0.1, 0, 100)
		n.Hygiene = clampF(n.Hygiene-0.05, 0, 100)
		if n.Energy >= 100 {
			s.Activity = "idle"
		}
		return
	}

	// auto sleep when out of energy
	if n.Energy <= 0 {
		if isNight {
			s.Activity = "sleeping"
		} else {
			s.Activity = "napping"
		}
		return
	}

	// decay rates per tick (~5s)
	n.Hunger = clampF(n.Hunger-decayRate(0.3, p.Patience), 0, 100)
	n.Energy = clampF(n.Energy-decayRate(0.2, p.EnergyLevel), 0, 100)
	n.Hygiene = clampF(n.Hygiene-decayRate(0.15, 50), 0, 100)
	n.Fun = clampF(n.Fun-decayRate(0.25, 50), 0, 100)

	// happiness drifts toward average of other needs
	avg := (n.Hunger + n.Energy + n.Hygiene + n.Fun + n.Health) / 5
	drift := (avg - n.Happiness) * 0.05
	n.Happiness = clampF(n.Happiness+drift, 0, 100)

	// health responds to sustained neglect
	low := 0
	for _, v := range []float64{n.Hunger, n.Energy, n.Hygiene} {
		if v < 20 {
			low++
		}
	}
	if low >= 2 {
		n.Health = clampF(n.Health-0.3, 0, 100)
	} else if n.Health < 100 && low == 0 {
		n.Health = clampF(n.Health+0.1, 0, 100)
	}

	// night drains energy faster
	if isNight {
		n.Energy = clampF(n.Energy-0.15, 0, 100)
	}
}

// SimulateOffline replays the ticks that would have happened while away,
// capped at 1 hour (720 ticks)
func (s *State) SimulateOffline() {
	if s.SavedAt == 0 {
		return
	}
	elapsed := time.Since(time.UnixMilli(s.SavedAt)).Seconds()
	ticks := int(elapsed / 5)
	if ticks > 720 {
		ticks = 720
	}
	for i := 0; i < ticks; i++ {
		s.Tick()
	}
}

// ==== actions ====

func (s *State) Feed(food string) {
	if s.IsAsleep() {
		return
	}
	n := &s.Needs
	p := s.Personality

	bonus := 1.0
	if food != "" && food == toString(s.FavoriteFood) {
		bonus = 1.5
	} else if food != "" && food == toString(s.WorstFood) {
		bonus = 0.3
	}
	picky := p.Pickiness / 100

	n.Hunger = clampF(n.Hunger+30*bonus, 0, 100)
	switch {
	case bonus > 1:
		n.Happiness = clampF(n.Happiness+10, 0, 100)
	case bonus < 1:
		n.Happiness = clampF(n.Happiness-5*picky, 0, 100)
	default:
		n.Happiness = clampF(n.Happiness+2, 0, 100)
	}
	s.addXP(5)
	s.logEvent("ate")
}

func (s *State) Play() {
	if s.IsAsleep() {
		return
	}
	n := &s.Needs
	n.Fun = clampF(n.Fun+25, 0, 100)
	n.Energy = clampF(n.Energy-0.5, 0, 100)
	n.Happiness = clampF(n.Happiness+10, 0, 100)
	s.addXP(10)
	s.logEvent("played")
}

func (s *State) Bathe() {
	if s.IsAsleep() {
		return
	}
	s.Needs.Hygiene = clampF(s.Needs.Hygiene+40, 0, 100)
	s.addXP(3)
	s.logEvent("bathed")
}

// Sleep toggles Mono into sleep, choosing nap or full sleep by time of day
func (s *State) Sleep() {
	if s.IsAsleep() {
		s.Activity = "idle"
		return
	}
	hour := time.Now().Hour()
	if hour >= 22 || hour < 7 {
		s.Activity = "sleeping"
	} else {
		s.Activity = "napping"
	}
}

func (s *State) SolveProblem() {
	s.CurrentProblem = nil
	s.ProblemsSolved++
	s.Needs.Happiness = clampF(s.Needs.Happiness+15, 0, 100)
	s.addXP(20)
	s.logEvent("solved a problem")
}

// ==== progress ====

// threshold = level * 50, XP rolls over on level-up
func (s *State) addXP(amount int) {
	s.XP += amount
	threshold := s.Level * 50
	if s.XP >= threshold {
		s.XP -= threshold
		s.Level++
		s.logEvent("leveled up")
	}
}

// XPThreshold returns XP needed for the next level
func (s *State) XPThreshold() int { return s.Level * 50 }

func (s *State) logEvent(desc string) {
	stamp := time.Now().Format("15:04") + "  " + desc
	s.EventHistory = append(s.EventHistory, stamp)
	if len(s.EventHistory) > 30 {
		s.EventHistory = s.EventHistory[len(s.EventHistory)-30:]
	}
}

func decayRate(base, personalityValue float64) float64 {
	return base * (1 + (50-personalityValue)/100)
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
