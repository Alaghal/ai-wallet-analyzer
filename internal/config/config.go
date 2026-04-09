package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv    string
	AppPort   int
	LogLevel  string
	LLMAPIURL string
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (Config, error) {
	port, err := getEnvAsInt("APP_PORT", 8080)
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	return Config{
		AppEnv:    getEnv("APP_ENV", "local"),
		AppPort:   port,
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LLMAPIURL: getEnv("LLM_API_URL", "http://localhost:11434"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	raw := getEnv(key, strconv.Itoa(fallback))
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer, got %q", key, raw)
	}
	return value, nil
}
