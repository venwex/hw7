package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"weather_api/internal/handler"
	"weather_api/internal/mocks"
	"weather_api/internal/models"
	"weather_api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupAuthHandler(mockAuth *mocks.MockAuthService) http.Handler {
	h := handler.NewAuthHandler(mockAuth, zap.NewNop())

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", h.Register)

	return mux
}

func TestRegisterHandler_Success(t *testing.T) {
	mockAuth := new(mocks.MockAuthService)
	router := setupAuthHandler(mockAuth)

	createdUser := models.User{
		ID:    1,
		Name:  "Alnur",
		Email: "test@gmail.com",
		Role:  "user",
	}

	mockAuth.
		On("Register", mock.Anything, "Alnur", "test@gmail.com", "password123").
		Return(createdUser, nil)

	body := `{"name":"Alnur","email":"test@gmail.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	assert.Contains(t, rec.Body.String(), `"id":1`)
	assert.Contains(t, rec.Body.String(), `"name":"Alnur"`)
	assert.Contains(t, rec.Body.String(), `"email":"test@gmail.com"`)
	assert.Contains(t, rec.Body.String(), `"role":"user"`)

	assert.NotContains(t, rec.Body.String(), "password")
	assert.NotContains(t, rec.Body.String(), "password_hash")

	mockAuth.AssertExpectations(t)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mockAuth := new(mocks.MockAuthService)
	router := setupAuthHandler(mockAuth)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid json")

	mockAuth.AssertNotCalled(t, "Register")
}

func TestRegisterHandler_InvalidInput(t *testing.T) {
	mockAuth := new(mocks.MockAuthService)
	router := setupAuthHandler(mockAuth)

	mockAuth.
		On("Register", mock.Anything, "", "test@gmail.com", "password123").
		Return(models.User{}, models.ErrInvalidInput)

	body := `{"name":"","email":"test@gmail.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid input")

	mockAuth.AssertExpectations(t)
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	mockAuth := new(mocks.MockAuthService)
	router := setupAuthHandler(mockAuth)

	mockAuth.
		On("Register", mock.Anything, "Alnur", "test@gmail.com", "password123").
		Return(models.User{}, service.ErrUserAlreadyExists)

	body := `{"name":"Alnur","email":"test@gmail.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "user already exists")

	mockAuth.AssertExpectations(t)
}

func TestRegisterHandler_InternalError(t *testing.T) {
	mockAuth := new(mocks.MockAuthService)
	router := setupAuthHandler(mockAuth)

	mockAuth.
		On("Register", mock.Anything, "Alnur", "test@gmail.com", "password123").
		Return(models.User{}, errors.New("database error"))

	body := `{"name":"Alnur","email":"test@gmail.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal server error")

	mockAuth.AssertExpectations(t)
}
