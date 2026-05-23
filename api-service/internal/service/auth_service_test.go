package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"weather_api/internal/auth"
	"weather_api/internal/mocks"
	"weather_api/internal/models"
	"weather_api/internal/repository"
	"weather_api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestAuthService() (*mocks.MockUserRepository, *service.AuthService) {
	tokenManager := auth.NewTokenManager("test-secret", 30*time.Minute)

	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, tokenManager)

	return mockRepo, authService
}

func TestRegister_Success(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	mockRepo.
		On("GetUserByEmail", mock.Anything, "test@gmail.com").
		Return(models.User{}, repository.ErrNotFound)

	createdUser := models.User{
		ID:    1,
		Name:  "Alnur",
		Email: "test@gmail.com",
		Role:  "user",
	}

	mockRepo.
		On("CreateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
			return user.Name == "Alnur" &&
				user.Email == "test@gmail.com" &&
				user.Role == "user" &&
				user.PasswordHash != "" &&
				user.PasswordHash != "password123"
		})).
		Return(createdUser, nil)

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"test@gmail.com",
		"password123",
	)

	require.NoError(t, err)

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Alnur", user.Name)
	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, "user", user.Role)

	mockRepo.AssertExpectations(t)
}

func TestRegister_EmptyName(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	user, err := authService.Register(
		context.Background(),
		"",
		"test@gmail.com",
		"password123",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, models.ErrInvalidInput)

	mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 0)
	mockRepo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestRegister_EmptyEmail(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"",
		"password123",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, models.ErrInvalidInput)

	mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 0)
	mockRepo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestRegister_EmptyPassword(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"test@gmail.com",
		"",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, models.ErrInvalidInput)

	mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 0)
	mockRepo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	existingUser := models.User{
		ID:    1,
		Name:  "Existing User",
		Email: "test@gmail.com",
		Role:  "user",
	}

	mockRepo.
		On("GetUserByEmail", mock.Anything, "test@gmail.com").
		Return(existingUser, nil)

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"test@gmail.com",
		"password123",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, service.ErrUserAlreadyExists)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestRegister_GetUserByEmailRepositoryError(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	dbErr := errors.New("database error")

	mockRepo.
		On("GetUserByEmail", mock.Anything, "test@gmail.com").
		Return(models.User{}, dbErr)

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"test@gmail.com",
		"password123",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, dbErr)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestRegister_CreateUserRepositoryError(t *testing.T) {
	mockRepo, authService := newTestAuthService()

	dbErr := errors.New("failed to create user")

	mockRepo.
		On("GetUserByEmail", mock.Anything, "test@gmail.com").
		Return(models.User{}, repository.ErrNotFound)

	mockRepo.
		On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).
		Return(models.User{}, dbErr)

	user, err := authService.Register(
		context.Background(),
		"Alnur",
		"test@gmail.com",
		"password123",
	)

	require.Error(t, err)

	assert.Equal(t, models.User{}, user)
	assert.ErrorIs(t, err, dbErr)

	mockRepo.AssertExpectations(t)
}
