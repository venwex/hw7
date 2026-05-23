package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"weather_api/internal/models"
	"weather_api/internal/repository"
)

type CityService struct {
	users  repository.UserRepository
	cities repository.CityRepository
}

func NewCityService(users repository.UserRepository, cities repository.CityRepository) *CityService {
	return &CityService{
		users:  users,
		cities: cities,
	}
}

/*
POST /users/{id}/cities — добавить город
GET /users/{id}/cities — список городов пользователя
DELETE /users/{id}/cities/{city_id} — удалить город

city string, userID
AddCity(ctx context.Context, city models.UserCity) (models.UserCity, error)
GetCities(ctx context.Context, userID int) ([]models.UserCity, error)
DeleteCity(ctx context.Context, userID, cityID int) error
*/

func (s *CityService) AddCity(ctx context.Context, userID int, city string) (models.UserCity, error) {
	if userID <= 0 {
		return models.UserCity{}, errors.New("invalid user id")
	}

	city = strings.TrimSpace(city)
	if city == "" {
		return models.UserCity{}, errors.New("city is required")
	}

	_, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserCity{}, fmt.Errorf("user not found")
	}

	return s.cities.AddCity(ctx, userID, city)
}

func (s *CityService) GetCities(ctx context.Context, userID int) ([]models.UserCity, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	_, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.cities.GetCities(ctx, userID)
}

func (s *CityService) DeleteCity(ctx context.Context, userID, cityID int) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}

	if cityID <= 0 {
		return errors.New("invalid city id")
	}

	_, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.cities.DeleteCity(ctx, userID, cityID)
}
