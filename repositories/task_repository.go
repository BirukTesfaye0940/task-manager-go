package repositories

import (
	"task-manager-go/models"
)

type TaskRepository interface {
	GetAll() ([]models.Task, error)
	GetByID(id int) (models.Task, error)
	Create(task models.Task) (models.Task, error)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
