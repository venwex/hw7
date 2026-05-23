package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"weather_gateway/internal/models"
)

type WeatherClient interface {
	GetWeather(ctx context.Context, city string) (models.WeatherResult, error)
}

type OpenMeteoClient struct {
	httpClient *http.Client
}

func NewOpenMeteoClient(httpClient *http.Client) *OpenMeteoClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &OpenMeteoClient{httpClient: httpClient}
}

type geocodingResponse struct {
	Results []geocodingResult `json:"results"`
}

type geocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

type forecastResponse struct {
	Current currentWeather `json:"current"`
}

type currentWeather struct {
	Temperature float64 `json:"temperature_2m"`
	WeatherCode int     `json:"weather_code"`
}

func (c *OpenMeteoClient) GetWeather(ctx context.Context, city string) (models.WeatherResult, error) {
	location, err := c.geocodeCity(ctx, city)
	if err != nil {
		return models.WeatherResult{}, err
	}

	weather, err := c.getCurrentWeather(ctx, location.Latitude, location.Longitude)
	if err != nil {
		return models.WeatherResult{}, err
	}

	return models.WeatherResult{
		City:        city,
		Temperature: weather.Temperature,
		Description: weatherCodeToDescription(weather.WeatherCode),
	}, nil
}

func (c *OpenMeteoClient) geocodeCity(ctx context.Context, city string) (geocodingResult, error) {
	values := url.Values{}
	values.Set("name", city)
	values.Set("count", "1")
	values.Set("language", "en")
	values.Set("format", "json")

	endpoint := "https://geocoding-api.open-meteo.com/v1/search?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return geocodingResult{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return geocodingResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return geocodingResult{}, fmt.Errorf("geocoding api returned status %d", resp.StatusCode)
	}

	var data geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return geocodingResult{}, err
	}

	if len(data.Results) == 0 {
		return geocodingResult{}, errors.New("city not found")
	}

	return data.Results[0], nil
}

func (c *OpenMeteoClient) getCurrentWeather(ctx context.Context, latitude, longitude float64) (currentWeather, error) {
	values := url.Values{}
	values.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	values.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	values.Set("current", "temperature_2m,weather_code")

	endpoint := "https://api.open-meteo.com/v1/forecast?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return currentWeather{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return currentWeather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return currentWeather{}, fmt.Errorf("forecast api returned status %d", resp.StatusCode)
	}

	var data forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return currentWeather{}, err
	}

	return data.Current, nil
}

func weatherCodeToDescription(code int) string {
	switch code {
	case 0:
		return "clear sky"
	case 1, 2, 3:
		return "partly cloudy"
	case 45, 48:
		return "foggy"
	case 51, 53, 55:
		return "drizzle"
	case 56, 57:
		return "freezing drizzle"
	case 61, 63, 65:
		return "rainy"
	case 66, 67:
		return "freezing rain"
	case 71, 73, 75:
		return "snowy"
	case 77:
		return "snow grains"
	case 80, 81, 82:
		return "rain showers"
	case 95:
		return "thunderstorm"
	case 96, 99:
		return "thunderstorm with hail"
	default:
		return "unknown"
	}
}
