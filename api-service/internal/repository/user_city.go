package repository

import (
	"context"
	"weather_api/internal/models"

	"github.com/jmoiron/sqlx"
)

type CityRepo struct {
	db *sqlx.DB
}

func NewCityRepo(db *sqlx.DB) *CityRepo {
	return &CityRepo{
		db: db,
	}
}

func (repo *CityRepo) AddCity(ctx context.Context, userID int, city string) (models.UserCity, error) {
	query := `
		INSERT INTO user_cities (user_id, city)
		VALUES ($1, $2)
		RETURNING id, user_id, city, created_at;
	`

	var created models.UserCity
	if err := repo.db.GetContext(ctx, &created, query, userID, city); err != nil {
		return models.UserCity{}, err
	}

	return created, nil
}

func (repo *CityRepo) GetCities(ctx context.Context, userID int) ([]models.UserCity, error) {
	query := `
		SELECT id, user_id, city, created_at
		FROM user_cities
		WHERE user_id = $1
		ORDER BY id;
	`

	var cities []models.UserCity
	if err := repo.db.SelectContext(ctx, &cities, query, userID); err != nil {
		return nil, err
	}

	return cities, nil
}

func (repo *CityRepo) DeleteCity(ctx context.Context, userID, cityID int) error {
	query := `
		DELETE FROM user_cities
		WHERE id = $1 AND user_id = $2;
	`

	_, err := repo.db.ExecContext(ctx, query, cityID, userID)
	return err
}
