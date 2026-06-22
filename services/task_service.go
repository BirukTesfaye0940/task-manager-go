package services

import (
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

func (s *TaskService) GetAllTasks() []models.Task {
	return s.repo.GetAll()
}

func (s *TaskService) GetTask(id int) (models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) CreateTask(task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, ErrTitleRequired
	}
	return s.repo.Create(task)
}

func (s *TaskService) UpdateTask(id int, task models.Task) (models.Task, error) {
	return s.repo.Update(id, task)
}

func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}
