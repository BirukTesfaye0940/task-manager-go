package services

import (
	"context"
	"testing"

	"task-manager-go/models"
	"task-manager-go/repositories"
)

type mockUserRepository struct {
	users map[string]models.User
}

func (m *mockUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	if _, exists := m.users[user.Username]; exists {
		return models.User{}, repositories.ErrUsernameTaken
	}
	user.ID = len(m.users) + 1
	m.users[user.Username] = user
	return user, nil
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (models.User, error) {
	user, exists := m.users[username]
	if !exists {
		return models.User{}, repositories.ErrUserNotFound
	}
	return user, nil
}

func setupAuthService() *AuthService {
	userRepo := &mockUserRepository{users: make(map[string]models.User)}
	return NewAuthService(userRepo, "test-secret-key")
}

func TestAuthService_Register(t *testing.T) {
	svc := setupAuthService()
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		user, err := svc.Register(ctx, "alice", "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.ID == 0 {
			t.Error("expected non-zero user ID")
		}
		if user.Username != "alice" {
			t.Errorf("expected username 'alice', got '%s'", user.Username)
		}
		if user.PasswordHash == "password123" {
			t.Error("password should be hashed, not stored as plaintext")
		}
	})

	t.Run("duplicate username returns error", func(t *testing.T) {
		_, err := svc.Register(ctx, "alice", "otherpassword")
		if err == nil {
			t.Fatal("expected error for duplicate username, got nil")
		}
	})

	t.Run("empty username returns error", func(t *testing.T) {
		_, err := svc.Register(ctx, "", "password")
		if err != ErrUsernameRequired {
			t.Errorf("expected ErrUsernameRequired, got %v", err)
		}
	})

	t.Run("empty password returns error", func(t *testing.T) {
		_, err := svc.Register(ctx, "bob", "")
		if err != ErrPasswordRequired {
			t.Errorf("expected ErrPasswordRequired, got %v", err)
		}
	})

	t.Run("canceled context returns error", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := svc.Register(cancelCtx, "carol", "password")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	svc := setupAuthService()
	ctx := context.Background()

	// Pre-register a user
	_, err := svc.Register(ctx, "testuser", "correctpassword")
	if err != nil {
		t.Fatalf("setup: register failed: %v", err)
	}

	t.Run("successful login returns token", func(t *testing.T) {
		token, err := svc.Login(ctx, "testuser", "correctpassword")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})

	t.Run("wrong password returns ErrInvalidCredentials", func(t *testing.T) {
		_, err := svc.Login(ctx, "testuser", "wrongpassword")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("non-existent user returns ErrInvalidCredentials", func(t *testing.T) {
		_, err := svc.Login(ctx, "nobody", "password")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("empty credentials return ErrInvalidCredentials", func(t *testing.T) {
		_, err := svc.Login(ctx, "", "")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("canceled context returns error", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := svc.Login(cancelCtx, "testuser", "correctpassword")
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}
