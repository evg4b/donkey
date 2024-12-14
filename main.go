package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evg4b/donkey/internal/config"
	"github.com/evg4b/donkey/internal/donkey"
	"github.com/evg4b/donkey/internal/store"
)

var Version = "X.X.X"

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error creating stire:", err)
		os.Exit(1)
	}

	version := false
	promt := ""
	pattern := ""

	flag.BoolVar(&version, "version", false, "Print version of the application")

	flag.StringVar(&config.DefaultProvider, "provider", config.DefaultProvider, "Default provider")
	flag.StringVar(&config.DefaultModel, "model", config.DefaultModel, "Default model")
	flag.DurationVar(&config.Timeout, "timeout", config.Timeout, "Timeout for request")

	flag.StringVar(&pattern, "pattern", pattern, "Pattern for the files")
	flag.StringVar(&promt, "promt", promt, "Promt for the model")

	flag.Parse()

	if version {
		fmt.Println("🫏 version:", Version)
		os.Exit(0)
	}

	store, err := store.NewStore(
		config.DefaultProvider,
		config.DefaultModel,
		config.Timeout,
	)

	if err != nil {
		fmt.Println("Error creating stire:", err)
		os.Exit(1)
	}

	model := donkey.InitialModel(store, pattern, promt)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
