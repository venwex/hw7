package service

import (
	"weather_api/internal/auth"
	"weather_api/internal/repository"
)
import "weather_api/internal/client"

type Service struct {
	User    *UserService
	City    *CityService
	Weather *WeatherService
	Auth    *AuthService
}

func NewService(repos *repository.Repository, weatherClient client.WeatherClient, tokenManager *auth.TokenManager) *Service {
	return &Service{
		User:    NewUserService(repos.Users),
		City:    NewCityService(repos.Users, repos.Cities),
		Weather: NewWeatherService(repos.Users, repos.Cities, repos.WeatherHistory, weatherClient),
		Auth:    NewAuthService(repos.Users, tokenManager),
	}
}
