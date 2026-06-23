package config

import (
	"errors"
	"os"
)

type Config struct {
	Addr     string
	DBPath   string
	APIToken string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:   getenv("KENDO_ADDR", ":8080"),
		DBPath: getenv("KENDO_DB_PATH", "/data/kendo.sqlite3"),
	}

	cfg.APIToken = os.Getenv("KENDO_API_TOKEN")
	if cfg.APIToken == "" {
		return Config{}, errors.New("KENDO_API_TOKEN is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
