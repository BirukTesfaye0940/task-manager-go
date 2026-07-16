package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"task-manager-go/config"
)

// Connect establishes a connection to the PostgreSQL database, configures
// connection pool tuning parameters, and verifies connectivity.
func Connect(ctx context.Context, cfg *config.Config) (*sql.DB, error) {

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
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Limit max concurrent active connections
	db.SetMaxIdleConns(25)                 // Match MaxOpenConns to avoid constant open/close churn
	db.SetConnMaxLifetime(5 * time.Minute) // Periodically recycle connections
	db.SetConnMaxIdleTime(5 * time.Minute) // Close idle connections that sit unused

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully (connection pool tuned)")
	return db, nil
}
