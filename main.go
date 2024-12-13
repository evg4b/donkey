package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evg4b/donkey/internal/donkey"
	"github.com/evg4b/donkey/internal/store"
)

func main() {
	store, err := store.NewStore()
	if err != nil {
		fmt.Println("Error creating stire:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(donkey.InitialModel(store)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
