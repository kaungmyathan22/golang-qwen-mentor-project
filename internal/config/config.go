package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration sourced from environment variables.
type Config struct {
	// AppEnv is required. Accepted values: development, staging, production.
	AppEnv string
	// Port is the TCP port the HTTP server listens on. Defaults to "8080".
	Port string
	// LogLevel controls zap verbosity. Defaults to "info".
	LogLevel string
}

// Load reads configuration from environment variables and validates required fields.
// It first attempts to load a .env file from the working directory; if the file
// does not exist the call is silently skipped (real env vars always take precedence).
func Load() (*Config, error) {
	// godotenv.Load only sets a key when it is NOT already present in the environment,
	// so actual env vars always win over the .env file.
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:   os.Getenv("APP_ENV"),
		Port:     envOr("APP_PORT", "8080"),
		LogLevel: envOr("LOG_LEVEL", "info"),
	}

	var errs []string

	if cfg.AppEnv == "" {
		errs = append(errs, "APP_ENV is required")
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Sprintf("APP_PORT %q is invalid", cfg.Port))
	}

	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[cfg.AppEnv] {
		errs = append(errs, fmt.Sprintf("APP_ENV %q is invalid", cfg.AppEnv))
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[strings.ToLower(cfg.LogLevel)] {
		errs = append(errs, fmt.Sprintf("LOG_LEVEL %q is invalid; must be one of debug|info|warn|error", cfg.LogLevel))
	} else {
		cfg.LogLevel = strings.ToLower(cfg.LogLevel)
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
