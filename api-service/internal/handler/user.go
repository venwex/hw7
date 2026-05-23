package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weather_api/internal/auth"
	"weather_api/internal/handler/dto"
	"weather_api/internal/models"
	"weather_api/internal/repository"
	u "weather_api/internal/utils"
)

type UserUseCase interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, id int, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id int) error
	Me(ctx context.Context, userID int) (models.User, error)
}

type UserHandler struct {
	users UserUseCase
}

func NewUserHandler(users UserUseCase) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.users.GetUsers(ctx)
	if err != nil {
		log.Printf("error getting users: %v", err)
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := u.GetID(r)
	if err != nil {
		log.Printf("error getting user id: %v", err)
		u.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.users.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			u.WriteError(w, http.StatusNotFound, "user not found")
			return

		case errors.Is(err, models.ErrInvalidID):
			u.WriteError(w, http.StatusBadRequest, "invalid id")
			return

		default:
			u.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	resp := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	u.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("error decoding user: %v", err)
		u.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.users.CreateUser(ctx, user)
	if err != nil {
		log.Printf("error creating user: %v", err)
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("error decoding user: %v", err)
		u.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := u.GetID(r)
	if err != nil {
		log.Printf("error getting user id: %v", err)
		u.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err = h.users.UpdateUser(ctx, id, user)
	if err != nil {
		log.Printf("error updating user: %v", err)
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) { // soft Delete
	ctx := r.Context()

	id, err := u.GetID(r)
	if err != nil {
		log.Printf("error getting user id: %v", err)
		u.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.users.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("error deleting user: %v", err)
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	currentUser, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.users.Me(r.Context(), currentUser.ID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	u.WriteJSON(w, http.StatusOK, user)
}
