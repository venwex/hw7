package main

import (
	"log/slog"
	"net/http"
	"os"

	"weather_gateway/internal/client"
	"weather_gateway/internal/handler"
	"weather_gateway/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	weatherClient := client.NewOpenMeteoClient(nil)
	weatherHandler := handler.NewWeatherHandler(weatherClient, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"gateway-service"}`))
	})
	mux.HandleFunc("GET /weather", weatherHandler.GetWeather)

	loggedMux := middleware.LoggingMiddleware(logger)(mux)

	addr := ":" + getEnv("PORT", "8081")
	logger.Info("gateway-service started", "addr", addr)
	if err := http.ListenAndServe(addr, loggedMux); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
