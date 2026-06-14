# hitori

A little life, in your terminal. Create a Mono, shape its personality, and watch it
live: it gets hungry, tired, bored, and happy on its own, and keeps living while
you're away.

```
   .---.
  ( ^ ^ )
  |  ◡  |
   '---'
   /| |\
```

## Setup

```sh
git clone https://github.com/mathiiiiiis/hitori-tui
cd hitori-tui
go mod tidy
go build -o hitori .
cp hitori ~/.local/bin/
```

## Usage

```sh
hitori           # launch the game (Discord login on first run)
hitori logout    # clear stored login
hitori version
```

First launch: log in via Discord, then a creation wizard walks you through naming
your Mono, setting a birthday and catchphrase, and shaping five personality traits.
After that you land in the world view.

### Screens

- **Create**: name, birthday, catchphrase, and personality sliders
- **World**: your Mono living its life: expression, mood, activity, all six needs, level/XP, recent events
- **Interact** (`i`): feed, play, bathe, sleep/wake, solve a problem
- **Customize** (`c`): re-tune personality traits

### Keys

| Key      | Where    | Action            |
|----------|----------|-------------------|
| `i`      | world    | open interactions |
| `c`      | world    | open customize    |
| `s`      | world    | sync now          |
| `↑/↓`    | menus    | select            |
| `←/→`    | sliders  | adjust            |
| `enter`  | menus    | confirm / do      |
| `esc`    | sub-view | back              |
| `q`      | world    | quit (saves)      |

## How a Mono works

Every 5 seconds Mono ticks: needs decay, happiness drifts toward the average of the
others, and health responds to sustained neglect. Personality changes the rates:
high Patience slows hunger, high Energy slows tiredness, high Optimism lifts mood.
Mood is derived from needs (sick / exhausted / tired / starving / bored / sad /
neutral / content / happy) and drives Mono's expression. When energy hits zero Mono
sleeps on its own (a nap by day, full sleep at night).

When you reopen the app, the time you were away is simulated forward (capped at one
hour) so Mono's state reflects real elapsed time.

## Backend

Lives at [hitori-backend](https://github.com/mathiiiiiis/hitori-backend). Handles
Discord/Google OAuth (including the CLI device-code flow this app uses) and the
`/save` endpoint that persists your Mono.
