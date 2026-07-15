package repositories

import (
	"context"
	"errors"
	"task-manager-go/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameTaken     = errors.New("username is already taken")
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
}
