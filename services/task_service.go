package services

import (
	"task-manager-go/models"
	"task-manager-go/repositories"
)

type TaskService struct {
	repo repositories.TaskRepository
}

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

func (s *TaskService) CreateTask(task models.Task) models.Task {
	return s.repo.Create(task)
}
