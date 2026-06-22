package repositories

import (
	"task-manager-go/models"
)

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

func (r *InMemoryTaskRepo) GetByID(id int) (models.Task, error) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return models.Task{}, ErrTaskNotFound
}

func (r *InMemoryTaskRepo) Create(task models.Task) (models.Task, error) {

	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *InMemoryTaskRepo) Update(id int, task models.Task) (models.Task, error) {

	for i, existingTask := range r.tasks {

		if existingTask.ID == id {

			task.ID = id

			r.tasks[i] = task

			return task, nil
		}
	}

	return models.Task{}, ErrTaskNotFound
}

func (r *InMemoryTaskRepo) Delete(id int) error {

	for i, task := range r.tasks {

		if task.ID == id {

			r.tasks = append(
				r.tasks[:i],
				r.tasks[i+1:]...,
			)

			return nil
		}
	}

	return ErrTaskNotFound
}
