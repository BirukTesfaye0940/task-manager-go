package services

import (
	"context"
	"testing"
	"task-manager-go/models"
	"task-manager-go/repositories"
)

func TestTaskService_ContextCanceled(t *testing.T) {
	repo := repositories.NewInMemoryTaskRepo()
	service := NewTaskService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel the context

	t.Run("GetAllTasks", func(t *testing.T) {
		_, err := service.GetAllTasks(ctx)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("GetTask", func(t *testing.T) {
		_, err := service.GetTask(ctx, 1)
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
		err := service.DeleteTask(ctx, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}
