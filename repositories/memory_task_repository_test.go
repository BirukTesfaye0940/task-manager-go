package repositories

import (
	"context"
	"testing"
	"task-manager-go/models"
)

func TestInMemoryTaskRepo_ContextCanceled(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel the context

	t.Run("GetAll", func(t *testing.T) {
		_, err := repo.GetAll(ctx)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Create", func(t *testing.T) {
		_, err := repo.Create(ctx, models.Task{Title: "Test"})
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		_, err := repo.Update(ctx, 1, models.Task{Title: "Test"})
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, 1)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}
