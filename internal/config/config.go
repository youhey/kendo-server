package config

import "os"

type Config struct {
	Addr   string
	DBPath string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:   getenv("KENDO_ADDR", ":8080"),
		DBPath: getenv("KENDO_DB_PATH", "/data/kendo.sqlite3"),
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
