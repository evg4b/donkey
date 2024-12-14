package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func WriteConfig(config *Config) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(homedir, fileName))
	if err != nil {
		return err
	}

	defer file.Close()

	return toml.NewEncoder(file).
		Encode(config)
}
