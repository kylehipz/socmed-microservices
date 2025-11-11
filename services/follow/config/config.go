package config

import "os"

type Config struct {
	DATABASE_URL string
}

func NewSettings() *Config {
	return &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}
}

var Settings = NewSettings()
