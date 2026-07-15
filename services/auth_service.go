package services

import (
	"context"
	"errors"
	"time"

	"task-manager-go/models"
	"task-manager-go/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameRequired   = errors.New("username is required")
	ErrPasswordRequired   = errors.New("password is required")
)

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret []byte
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (models.User, error) {
	if err := ctx.Err(); err != nil {
		return models.User{}, err
	}

	if username == "" {
		return models.User{}, ErrUsernameRequired
	}
	if password == "" {
		return models.User{}, ErrPasswordRequired
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hashedBytes),
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if username == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
