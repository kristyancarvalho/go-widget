package config

import (
	"os"
)

type Config struct {
	GitHubUsername string
	GitHubToken    string
	Port           string
}

func Load() *Config {
	return &Config{
		GitHubUsername: getEnv("GITHUB_USERNAME", ""),
		GitHubToken:    getEnv("GITHUB_TOKEN", ""),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
