package handler

import (
	"net/http"
	"weather_api/internal/auth"
	"weather_api/internal/service"
	u "weather_api/internal/utils"
)

type WeatherHandler struct {
	Weather *service.WeatherService
}

func NewWeatherHandler(weather *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		Weather: weather,
	}
}

func (h *WeatherHandler) GetUserWeather(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		u.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	weather, err := h.Weather.GetUserWeather(r.Context(), user.ID)
	if err != nil {
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusOK, weather)
}

func (h *WeatherHandler) GetWeatherHistory(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		u.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		u.WriteError(w, http.StatusBadRequest, "city is required")
		return
	}

	history, err := h.Weather.GetWeatherHistory(r.Context(), user.ID, city, 10)
	if err != nil {
		u.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.WriteJSON(w, http.StatusOK, history)
}
