package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather_api/internal/mocks"
	"weather_api/internal/models"
	"weather_api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserHandler(mockUsers *mocks.MockUserService) http.Handler {
	h := NewUserHandler(mockUsers)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", h.GetUserByID)

	return mux
}

func TestGetUserByIDHandler_Success(t *testing.T) {
	mockUsers := new(mocks.MockUserService)
	router := setupUserHandler(mockUsers)

	user := models.User{
		ID:    1,
		Name:  "Alnur",
		Email: "test@gmail.com",
		Role:  "user",
	}

	mockUsers.
		On("GetUserByID", mock.Anything, 1).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":1`)
	assert.Contains(t, rec.Body.String(), `"name":"Alnur"`)
	assert.Contains(t, rec.Body.String(), `"email":"test@gmail.com"`)
	assert.Contains(t, rec.Body.String(), `"role":"user"`)

	assert.NotContains(t, rec.Body.String(), "password")
	assert.NotContains(t, rec.Body.String(), "password_hash")

	mockUsers.AssertExpectations(t)
}

func TestGetUserByIDHandler_InvalidID(t *testing.T) {
	mockUsers := new(mocks.MockUserService)
	router := setupUserHandler(mockUsers)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockUsers.AssertNotCalled(t, "GetUserByID")
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	mockUsers := new(mocks.MockUserService)
	router := setupUserHandler(mockUsers)

	mockUsers.
		On("GetUserByID", mock.Anything, 999).
		Return(models.User{}, repository.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockUsers.AssertExpectations(t)
}

func TestGetUserByIDHandler_InternalError(t *testing.T) {
	mockUsers := new(mocks.MockUserService)
	router := setupUserHandler(mockUsers)

	mockUsers.
		On("GetUserByID", mock.Anything, 1).
		Return(models.User{}, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockUsers.AssertExpectations(t)
}
