package repositories

import (
	"context"
	"task-manager-go/models"
)

type TaskRepository interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, error)
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Update(ctx context.Context, id int, task models.Task) (models.Task, error)
	Delete(ctx context.Context, id int) error
}
