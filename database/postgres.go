package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"task-manager-go/config"
)

func Connect() (*sql.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL")
	return db, nil
}
