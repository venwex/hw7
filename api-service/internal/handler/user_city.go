package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"weather_api/internal/auth"
	"weather_api/internal/service"
	u "weather_api/internal/utils"
)

type CityHandler struct {
	Cities *service.CityService
}

func NewCityHandler(cities *service.CityService) *CityHandler {
	return &CityHandler{Cities: cities}
}

type addCityRequest struct {
	City string `json:"city"`
}

func (h *CityHandler) AddCity(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req addCityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	city, err := h.Cities.AddCity(r.Context(), user.ID, req.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.WriteJSON(w, http.StatusCreated, city)
}

func (h *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cities, err := h.Cities.GetCities(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.WriteJSON(w, http.StatusOK, cities)
}

func (h *CityHandler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cityID, err := strconv.Atoi(r.PathValue("city_id"))
	if err != nil || cityID <= 0 {
		http.Error(w, "invalid city id", http.StatusBadRequest)
		return
	}

	err = h.Cities.DeleteCity(r.Context(), user.ID, cityID)
	if err != nil {
		http.Error(w, "city not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
