package services

import (
	"context"
	"testing"
	"task-manager-go/models"
)

type mockTaskRepository struct{}

func (m mockTaskRepository) GetAll(ctx context.Context, userID int) ([]models.Task, error) {
	return nil, nil
}
func (m mockTaskRepository) GetByID(ctx context.Context, id int, userID int) (models.Task, error) {
	return models.Task{}, nil
}
func (m mockTaskRepository) Create(ctx context.Context, task models.Task) (models.Task, error) {
	return models.Task{}, nil
}
func (m mockTaskRepository) Update(ctx context.Context, id int, task models.Task) (models.Task, error) {
	return models.Task{}, nil
}
func (m mockTaskRepository) Delete(ctx context.Context, id int, userID int) error {
	return nil
}

func TestTaskService_ContextCanceled(t *testing.T) {
	repo := mockTaskRepository{}
	service := NewTaskService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel the context

	t.Run("GetAllTasks", func(t *testing.T) {
		_, err := service.GetAllTasks(ctx, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("GetTask", func(t *testing.T) {
		_, err := service.GetTask(ctx, 1, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("CreateTask", func(t *testing.T) {
		_, err := service.CreateTask(ctx, models.Task{Title: "Test"})
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("UpdateTask", func(t *testing.T) {
		_, err := service.UpdateTask(ctx, 1, models.Task{Title: "Test"})
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		err := service.DeleteTask(ctx, 1, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}
