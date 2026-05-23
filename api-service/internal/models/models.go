package models

import "time"

type User struct {
	ID           int        `db:"id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	Role         string     `db:"role"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type UserCity struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	City      string    `db:"city"`
	CreatedAt time.Time `db:"created_at"`
}

type WeatherResult struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type UserWeatherResponse struct {
	UserID  int             `json:"user_id"`
	Weather []WeatherResult `json:"weather"`
}

type WeatherHistory struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	City        string    `db:"city" json:"city"`
	Temperature float64   `db:"temperature" json:"temperature"`
	Description string    `db:"description" json:"description"`
	RequestedAt time.Time `db:"requested_at" json:"requested_at"`
}

type WeatherHistoryItem struct {
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}

type WeatherHistoryResponse struct {
	UserID  int                  `json:"user_id"`
	City    string               `json:"city"`
	History []WeatherHistoryItem `json:"history"`
}
