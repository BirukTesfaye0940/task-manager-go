package config

import (
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

func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-jwt-secret-key-change-in-prod"
	}

	return &Config{

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		ServerPort: os.Getenv("SERVER_PORT"),
		JWTSecret:  jwtSecret,
	}, nil
}
