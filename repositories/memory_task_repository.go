package repositories

import (
	"context"
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

func (r *InMemoryTaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.tasks, nil
}

func (r *InMemoryTaskRepo) GetByID(ctx context.Context, id int) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return models.Task{}, ErrTaskNotFound
}

func (r *InMemoryTaskRepo) Create(ctx context.Context, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}

	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *InMemoryTaskRepo) Update(ctx context.Context, id int, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}

	for i, existingTask := range r.tasks {

		if existingTask.ID == id {

			task.ID = id

			r.tasks[i] = task

			return task, nil
		}
	}

	return models.Task{}, ErrTaskNotFound
}

func (r *InMemoryTaskRepo) Delete(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

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
