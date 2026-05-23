package handler

import (
	"weather_api/internal/service"

	"go.uber.org/zap"
)

type Handler struct {
	Users   *UserHandler
	Cities  *CityHandler
	Weather *WeatherHandler
	Auth    *AuthHandler
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		Users:   NewUserHandler(svc.User),
		Cities:  NewCityHandler(svc.City),
		Weather: NewWeatherHandler(svc.Weather),
		Auth:    NewAuthHandler(svc.Auth, zap.NewNop()),
	}
}
