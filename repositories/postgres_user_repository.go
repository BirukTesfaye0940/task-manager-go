package repositories

import (
	"context"
	"database/sql"
	"task-manager-go/models"

	"github.com/lib/pq"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	if err := ctx.Err(); err != nil {
		return models.User{}, err
	}

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, created_at",
		user.Username, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique_violation code in postgres
			return models.User{}, ErrUsernameTaken
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	if err := ctx.Err(); err != nil {
		return models.User{}, err
	}

	var user models.User
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}
