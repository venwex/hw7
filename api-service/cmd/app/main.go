package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"weather_api/internal/config"

	"weather_api/internal/auth"
	"weather_api/internal/client"
	"weather_api/internal/handler"
	"weather_api/internal/middleware"
	"weather_api/internal/migrations"
	"weather_api/internal/repository"
	"weather_api/internal/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lpernett/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	_ = godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	ttlRaw := os.Getenv("ACCESS_TOKEN_TTL")
	if ttlRaw == "" {
		ttlRaw = "30m"
	}

	accessTTL, err := time.ParseDuration(ttlRaw)
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := auth.NewTokenManager(jwtSecret, accessTTL)

	cfg := config.NewPostgresConfig()

	db, err := sqlx.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Connection error to the DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	if err := migrations.Run(context.Background(), db, os.Getenv("MIGRATIONS_PATH")); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	repo := repository.NewRepository(db)

	gatewayURL := os.Getenv("GATEWAY_BASE_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8081"
	}
	weatherClient := client.NewGatewayWeatherClient(gatewayURL, nil)

	svc := service.NewService(repo, weatherClient, tokenManager)

	h := handler.NewHandler(svc)

	mux := initRoutes(h, tokenManager)
	loggedMux := middleware.LoggingMiddleware(logger)(mux)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}

func initRoutes(h *handler.Handler, tokenManager *auth.TokenManager) *http.ServeMux {
	mux := http.NewServeMux()

	authMW := middleware.AuthMiddleware(tokenManager)
	adminMW := middleware.RequireRole("admin")

	protected := func(fn http.HandlerFunc) http.Handler {
		return authMW(fn)
	}

	adminOnly := func(fn http.HandlerFunc) http.Handler {
		return authMW(adminMW(fn))
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"api-service"}`))
	})

	mux.HandleFunc("POST /auth/register", h.Auth.Register)
	mux.HandleFunc("POST /auth/login", h.Auth.Login)

	mux.Handle("GET /users/me", protected(h.Users.Me))

	mux.Handle("POST /cities", protected(h.Cities.AddCity))
	mux.Handle("GET /cities", protected(h.Cities.GetCities))
	mux.Handle("DELETE /cities/{city_id}", protected(h.Cities.DeleteCity))

	mux.Handle("GET /weather", protected(h.Weather.GetUserWeather))
	mux.Handle("GET /weather/history", protected(h.Weather.GetWeatherHistory))

	mux.Handle("GET /users", adminOnly(h.Users.GetUsers))
	mux.Handle("GET /users/{id}", adminOnly(h.Users.GetUserByID))
	mux.Handle("DELETE /users/{id}", adminOnly(h.Users.DeleteUser))

	return mux
}
