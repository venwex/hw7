package service

import (
	"context"
	m "weather_api/internal/models"
	"weather_api/internal/repository"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{
		users: users,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]m.User, error) {
	return s.users.GetUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (m.User, error) {
	if id <= 0 {
		return m.User{}, m.ErrInvalidID
	}

	return s.users.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user m.User) (m.User, error) {
	if len(user.Name) == 0 {
		return m.User{}, m.ErrInvalidName
	}

	if len(user.Email) == 0 {
		return m.User{}, m.ErrInvalidEmail
	}

	return s.users.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, user m.User) (m.User, error) {
	if id <= 0 {
		return m.User{}, m.ErrInvalidID
	}

	if len(user.Name) == 0 {
		return m.User{}, m.ErrInvalidName
	}

	if len(user.Email) == 0 {
		return m.User{}, m.ErrInvalidEmail
	}

	return s.users.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return m.ErrInvalidID
	}

	return s.users.DeleteUser(ctx, id)
}

func (s *UserService) Me(ctx context.Context, userID int) (m.User, error) {
	return s.users.GetUserByID(ctx, userID)
}
