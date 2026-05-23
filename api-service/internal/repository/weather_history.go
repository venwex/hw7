package repository

import (
	"context"
	"weather_api/internal/models"

	"github.com/jmoiron/sqlx"
)

type WeatherHistoryRepo struct {
	db *sqlx.DB
}

func NewWeatherHistoryRepository(db *sqlx.DB) *WeatherHistoryRepo {
	return &WeatherHistoryRepo{
		db: db,
	}
}

func (repo *WeatherHistoryRepo) CreateWeatherHistory(ctx context.Context, item models.WeatherHistory) (models.WeatherHistory, error) {
	query := `
		INSERT INTO weather_history (user_id, city, temperature, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, city, temperature, description, requested_at;
	`

	var created models.WeatherHistory
	if err := repo.db.GetContext(
		ctx,
		&created,
		query,
		item.UserID,
		item.City,
		item.Temperature,
		item.Description,
	); err != nil {
		return models.WeatherHistory{}, err
	}

	return created, nil
}

func (repo *WeatherHistoryRepo) GetWeatherHistoryByCity(ctx context.Context, userID int, city string, limit int) ([]models.WeatherHistory, error) {
	query := `
		SELECT id, user_id, city, temperature, description, requested_at
		FROM weather_history
		WHERE user_id = $1 AND city = $2
		ORDER BY requested_at DESC
		LIMIT $3;
	`

	var history []models.WeatherHistory
	if err := repo.db.SelectContext(ctx, &history, query, userID, city, limit); err != nil {
		return nil, err
	}

	return history, nil
}
