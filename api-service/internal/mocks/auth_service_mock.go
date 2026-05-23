package mocks

import (
	"context"
	"weather_api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, name, email, password string) (models.User, error) {
	args := m.Called(ctx, name, email, password)

	user, _ := args.Get(0).(models.User)
	return user, args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)

	token, _ := args.Get(0).(string)
	return token, args.Error(1)
}
