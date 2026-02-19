package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      int
	ProbePort int
	LogLevel  string
	DebugMode bool
	Database  DatabaseConfig
}

type DatabaseConfig struct {
	MigrationEnabled bool
	Host             string
	Port             int
	Name             string
	User             string
	Password         string
	SSLMode          string
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignored in production)
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnvInt("PORT", 8080),
		ProbePort: getEnvInt("PROBE_PORT", 8081),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		DebugMode: getEnvBool("DEBUG_MODE", false),
		Database: DatabaseConfig{
			MigrationEnabled: getEnvBool("MIGRATION_ENABLED", true),
			Host:             getEnv("DATABASE_HOST", "localhost"),
			Port:             getEnvInt("DATABASE_PORT", 5432),
			Name:             getEnv("DATABASE_NAME", "{{PROJECT}}"),
			User:             getEnv("DATABASE_USER", "{{PROJECT}}"),
			Password:         getEnv("DATABASE_PASSWORD", ""),
			SSLMode:          getEnv("DATABASE_SSLMODE", "disable"),
		},
	}

	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DATABASE_PASSWORD is required")
	}

	return cfg, nil
}

// ConnectionString returns the PostgreSQL connection string.
func (d DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host, d.Port, d.Name, d.SSLMode)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
