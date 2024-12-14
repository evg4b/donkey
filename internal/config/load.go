package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func LoadConfig() (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Join(homedir, fileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &DefaultConfig, WriteConfig(&DefaultConfig)
		}

		return nil, err
	}

	defer file.Close()

	decoder := toml.NewDecoder(file)
	var config Config
	if _, err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
