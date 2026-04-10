package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv          string
	AppPort         int
	LogLevel        string
	LLMAPIURL       string
	ProviderType    string
	EtherscanAPIURL string
	EtherscanAPIKey string
	HTTPTimeout     time.Duration
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

	timeoutSeconds, err := getEnvAsInt("HTTP_TIMEOUT_SECONDS", 5)
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	return Config{
		AppEnv:          getEnv("APP_ENV", "local"),
		AppPort:         port,
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		LLMAPIURL:       getEnv("LLM_API_URL", "http://localhost:11434"),
		ProviderType:    getEnv("PROVIDER_TYPE", "mock"),
		EtherscanAPIURL: getEnv("ETHERSCAN_API_URL", "https://api.etherscan.io/api"),
		EtherscanAPIKey: getEnv("ETHERSCAN_API_KEY", ""),
		HTTPTimeout:     time.Duration(timeoutSeconds) * time.Second,
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
