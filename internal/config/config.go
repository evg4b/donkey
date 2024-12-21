package config

import "time"

const fileName = ".donkey.toml"

type Config struct {
	DefaultProvider string
	DefaultModel    string
	Timeout         time.Duration
	DefaultSuffix   string
}

var DefaultConfig = Config{
	DefaultProvider: "ollama",
	DefaultModel:    "mistral-small:latest",
	Timeout:         time.Hour,
	DefaultSuffix:   "",
}
