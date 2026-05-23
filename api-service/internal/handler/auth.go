package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"weather_api/internal/models"
	u "weather_api/internal/utils"

	"weather_api/internal/handler/dto"
	"weather_api/internal/service"

	"go.uber.org/zap"
)

type AuthUseCase interface {
	Register(ctx context.Context, name, email, password string) (models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	auth   AuthUseCase
	logger *zap.Logger
}

func NewAuthHandler(auth AuthUseCase, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		auth:   auth,
		logger: logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid register json",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)

		u.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user, err := h.auth.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidInput):
			h.logger.Warn("invalid register input",
				zap.String("email", req.Email),
				zap.Error(err),
			)

			u.WriteError(w, http.StatusBadRequest, "invalid input")
			return

		case errors.Is(err, service.ErrUserAlreadyExists):
			h.logger.Warn("register failed: user already exists",
				zap.String("email", req.Email),
				zap.Error(err),
			)

			u.WriteError(w, http.StatusConflict, "user already exists")
			return

		default:
			h.logger.Error("register failed: internal error",
				zap.String("email", req.Email),
				zap.Error(err),
			)

			u.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	resp := dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	h.logger.Info("user registered successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("role", user.Role),
	)

	u.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid login json",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)

		u.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)

		u.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.logger.Info("user logged in successfully",
		zap.String("email", req.Email),
	)

	u.WriteJSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken: token,
	})
}
