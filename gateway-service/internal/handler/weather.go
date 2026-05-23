package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"weather_gateway/internal/client"
	"weather_gateway/internal/utils"
)

type WeatherHandler struct {
	weatherClient client.WeatherClient
	logger        *slog.Logger
}

func NewWeatherHandler(weatherClient client.WeatherClient, logger *slog.Logger) *WeatherHandler {
	return &WeatherHandler{weatherClient: weatherClient, logger: logger}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	if city == "" {
		utils.WriteError(w, http.StatusBadRequest, "city is required")
		return
	}

	result, err := h.weatherClient.GetWeather(r.Context(), city)
	if err != nil {
		h.logger.Error("get weather from external api", "city", city, "error", err)
		utils.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}
