package service

import (
	"context"
	"errors"
	"strings"

	"weather_api/internal/auth"
	"weather_api/internal/models"
	"weather_api/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type AuthService struct {
	users        repository.UserRepository
	tokenManager *auth.TokenManager
}

func NewAuthService(users repository.UserRepository, tokenManager *auth.TokenManager) *AuthService {
	return &AuthService{
		users:        users,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (models.User, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return models.User{}, models.ErrInvalidInput
	}

	_, err := s.users.GetUserByEmail(ctx, email)
	if err == nil {
		return models.User{}, ErrUserAlreadyExists
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return models.User{}, err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	return s.users.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
