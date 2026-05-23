package repository_test

import (
	"context"
	"os"
	"testing"
	"weather_api/internal/models"
	"weather_api/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5437/weather_test?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)

	schema := `
	DROP TABLE IF EXISTS users;

	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		role TEXT NOT NULL DEFAULT 'user',
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ
	);
	`

	_, err = db.Exec(schema)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.Exec("DROP TABLE IF EXISTS users;")
		_ = db.Close()
	})

	return db
}

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)

	repo := repository.NewUserRepo(db)

	ctx := context.Background()

	user := models.User{
		Name:         "Alnur",
		Email:        "alnur@test.com",
		Role:         "user",
		PasswordHash: "hashed-password",
	}

	createdUser, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, "Alnur", createdUser.Name)
	assert.Equal(t, "alnur@test.com", createdUser.Email)
	assert.Equal(t, "user", createdUser.Role)

	foundUser, err := repo.GetUserByID(ctx, createdUser.ID)
	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, "Alnur", foundUser.Name)
	assert.Equal(t, "alnur@test.com", foundUser.Email)
	assert.Equal(t, "user", foundUser.Role)
	assert.Equal(t, "hashed-password", foundUser.PasswordHash)
}
