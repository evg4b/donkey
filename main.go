package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evg4b/donkey/internal/donkey"
	"os"
)

func main() {
	if _, err := tea.NewProgram(donkey.InitialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
