package services

import (
	"context"
	"errors"
	"task-manager-go/models"
	"task-manager-go/repositories"
)

type TaskService struct {
	repo repositories.TaskRepository
}

var ErrTitleRequired = errors.New("title is required")

func NewTaskService(
	repo repositories.TaskRepository,
) *TaskService {

	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) GetAllTasks(ctx context.Context, userID int) ([]models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, userID)
}

func (s *TaskService) GetTask(ctx context.Context, id int, userID int) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	return s.repo.GetByID(ctx, id, userID)
}

func (s *TaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	if task.Title == "" {
		return models.Task{}, ErrTitleRequired
	}
	return s.repo.Create(ctx, task)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	return s.repo.Update(ctx, id, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int, userID int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id, userID)
}
