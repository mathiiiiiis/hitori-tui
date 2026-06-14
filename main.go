package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	hauth "github.com/mathiiiiiis/hitori/internal/auth"
	"github.com/mathiiiiiis/hitori/internal/tui"
)

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "logout":
			if err := hauth.DeleteToken(); err != nil {
				fmt.Fprintln(os.Stderr, "error clearing token:", err)
				os.Exit(1)
			}
			fmt.Println("logged out.")
			return
		case "version":
			fmt.Println("hitori v0.1.0")
			return
		case "help", "-h", "--help":
			printHelp()
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	}

	token := hauth.LoadToken()
	m := tui.New(token)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`hitori: a little life, in your terminal

USAGE
  hitori           launch the game
  hitori logout    clear stored login
  hitori version

On first run you'll log in via Discord, then create your Mono.
Your Mono syncs to the backend and keeps living while you're away.
`)
}
