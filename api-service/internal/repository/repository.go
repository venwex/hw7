package repository

import (
	"context"
	"weather_api/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, id int, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type CityRepository interface {
	AddCity(ctx context.Context, userID int, city string) (models.UserCity, error)
	GetCities(ctx context.Context, userID int) ([]models.UserCity, error)
	DeleteCity(ctx context.Context, userID, cityID int) error
}

type WeatherHistoryRepository interface {
	CreateWeatherHistory(ctx context.Context, item models.WeatherHistory) (models.WeatherHistory, error)
	GetWeatherHistoryByCity(ctx context.Context, userID int, city string, limit int) ([]models.WeatherHistory, error)
}

type Repository struct {
	Users          UserRepository
	Cities         CityRepository
	WeatherHistory WeatherHistoryRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users:          NewUserRepo(db),
		Cities:         NewCityRepo(db),
		WeatherHistory: NewWeatherHistoryRepository(db),
	}
}
