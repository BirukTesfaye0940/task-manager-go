package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	ServerPort string
	JWTSecret  string
}

// Load reads config from environment variables (loading from .env file if available).
// It returns an error if any of the required variables are missing.
func Load() (*Config, error) {
	// Attempt to load .env, but ignore error if the file doesn't exist
	// (common in containerized/production environments).
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration values are present.
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return errors.New("DB_HOST is required")
	}
	if c.DBPort == "" {
		return errors.New("DB_PORT is required")
	}
	if c.DBUser == "" {
		return errors.New("DB_USER is required")
	}
	if c.DBPassword == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.DBName == "" {
		return errors.New("DB_NAME is required")
	}
	if c.ServerPort == "" {
		c.ServerPort = "8080" // Set default port if not provided
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}
