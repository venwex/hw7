package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"weather_api/internal/models"
)

type WeatherClient interface {
	GetWeather(ctx context.Context, city string) (models.WeatherResult, error)
}

type GatewayWeatherClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGatewayWeatherClient(baseURL string, httpClient *http.Client) *GatewayWeatherClient {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://localhost:8081"
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	return &GatewayWeatherClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *GatewayWeatherClient) GetWeather(ctx context.Context, city string) (models.WeatherResult, error) {
	endpoint := c.baseURL + "/weather?city=" + url.QueryEscape(city)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return models.WeatherResult{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.WeatherResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return models.WeatherResult{}, fmt.Errorf("gateway returned status %d", resp.StatusCode)
	}

	var result models.WeatherResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.WeatherResult{}, err
	}

	return result, nil
}
