package service

import (
	"context"
	"errors"
	"strings"
	"weather_api/internal/client"
	"weather_api/internal/models"
	"weather_api/internal/repository"
)

type WeatherService struct {
	users          repository.UserRepository
	cities         repository.CityRepository
	weatherHistory repository.WeatherHistoryRepository
	weatherClient  client.WeatherClient
}

func NewWeatherService(
	users repository.UserRepository,
	cities repository.CityRepository,
	weatherHistory repository.WeatherHistoryRepository,
	weatherClient client.WeatherClient,
) *WeatherService {
	return &WeatherService{
		users:          users,
		cities:         cities,
		weatherHistory: weatherHistory,
		weatherClient:  weatherClient,
	}
}

func (s *WeatherService) GetUserWeather(ctx context.Context, userID int) (models.UserWeatherResponse, error) {
	if userID <= 0 {
		return models.UserWeatherResponse{}, errors.New("invalid user id")
	}

	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserWeatherResponse{}, err
	}

	userCities, err := s.cities.GetCities(ctx, user.ID)
	if err != nil {
		return models.UserWeatherResponse{}, err
	}

	response := models.UserWeatherResponse{
		UserID:  user.ID,
		Weather: make([]models.WeatherResult, 0, len(userCities)),
	}

	for _, userCity := range userCities {
		weather, err := s.weatherClient.GetWeather(ctx, userCity.City)
		if err != nil {
			return models.UserWeatherResponse{}, err
		}

		response.Weather = append(response.Weather, weather)

		_, err = s.weatherHistory.CreateWeatherHistory(ctx, models.WeatherHistory{
			UserID:      user.ID,
			City:        userCity.City,
			Temperature: weather.Temperature,
			Description: weather.Description,
		})
		if err != nil {
			return models.UserWeatherResponse{}, err
		}
	}

	return response, nil
}

var ErrCityRequired = errors.New("city is required")

func (s *WeatherService) GetWeatherHistory(ctx context.Context, userID int, city string, limit int) (models.WeatherHistoryResponse, error) {
	if userID <= 0 {
		return models.WeatherHistoryResponse{}, errors.New("invalid user id")
	}

	city = strings.TrimSpace(city)
	if city == "" {
		return models.WeatherHistoryResponse{}, ErrCityRequired
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return models.WeatherHistoryResponse{}, err
	}

	history, err := s.weatherHistory.GetWeatherHistoryByCity(ctx, user.ID, city, limit)
	if err != nil {
		return models.WeatherHistoryResponse{}, err
	}

	items := make([]models.WeatherHistoryItem, 0, len(history))
	for _, h := range history {
		items = append(items, models.WeatherHistoryItem{
			Temperature: h.Temperature,
			Description: h.Description,
			RequestedAt: h.RequestedAt,
		})
	}

	return models.WeatherHistoryResponse{
		UserID:  user.ID,
		City:    city,
		History: items,
	}, nil
}
