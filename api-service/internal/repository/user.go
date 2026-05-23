package repository

import (
	"context"
	"database/sql"
	"errors"

	"weather_api/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id;
	`

	var users []models.User
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at, deleted_at, role, password_hash
		FROM users
		WHERE id = $1 AND deleted_at IS NULL;
	`

	var user models.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, role, password_hash, created_at, updated_at, deleted_at
	`

	if user.Role == "" {
		user.Role = "user"
	}

	var created models.User
	err := r.db.GetContext(
		ctx,
		&created,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	return created, err
}

func (r *UserRepo) UpdateUser(ctx context.Context, id int, user models.User) (models.User, error) {
	query := `
		UPDATE users
		SET name = $1,
		    email = $2,
		    updated_at = now()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, name, email, created_at, updated_at, deleted_at;
	`

	var updated models.User
	if err := r.db.GetContext(ctx, &updated, query, user.Name, user.Email, id); err != nil {
		return models.User{}, err
	}

	return updated, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int) error {
	query := `
		UPDATE users
		SET deleted_at = now(),
		    updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}

	return user, nil
}
