package repositories

import (
	"errors"
	"task-manager-go/models"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	GetAll() []models.Task
	GetByID(id int) (models.Task, error)
	Create(task models.Task) models.Task
}

type InMemoryTaskRepo struct {
	tasks  []models.Task
	nextID int
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks:  []models.Task{},
		nextID: 1,
	}
}

func (r *InMemoryTaskRepo) GetAll() []models.Task {
	return r.tasks
}

func (r *InMemoryTaskRepo) GetbyID(id int) (models.Task, error) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return models.Task{}, ErrTaskNotFound
}

func (r *InMemoryTaskRepo) Create(task models.Task) models.Task {

	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task
}
