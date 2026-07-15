package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"task-manager-go/config"
)

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
		return nil, err
	}

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	if err := initSchema(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	fmt.Println("Connected to PostgreSQL and schema initialized")
	return db, nil
}

func initSchema(ctx context.Context, db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		done BOOLEAN DEFAULT FALSE
	);

	ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE;
	`
	_, err := db.ExecContext(ctx, schema)
	return err
}
