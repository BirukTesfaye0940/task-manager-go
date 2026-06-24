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

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetTask(ctx context.Context, id int) (models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, ErrTitleRequired
	}
	return s.repo.Create(ctx, task)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, task models.Task) (models.Task, error) {
	return s.repo.Update(ctx, id, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
