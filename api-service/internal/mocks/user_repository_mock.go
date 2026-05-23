package mocks

import (
	"context"
	"weather_api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)

	users, _ := args.Get(0).([]models.User)
	return users, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	args := m.Called(ctx, id)

	user, _ := args.Get(0).(models.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)

	createdUser, _ := args.Get(0).(models.User)
	return createdUser, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id int, user models.User) (models.User, error) {
	args := m.Called(ctx, id, user)

	updatedUser, _ := args.Get(0).(models.User)
	return updatedUser, args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)

	user, _ := args.Get(0).(models.User)
	return user, args.Error(1)
}
