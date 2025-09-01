package open_weather

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/conf"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchWeatherByLocation(t *testing.T) {
	tests := []struct {
		name           string
		cityName       string
		country        string
		apiKey         string
		mockResponse   Response
		mockStatusCode int
		expectedResult *schemata.FetchWeatherResponse
		expectedError  error
	}{
		{
			name:     "successful weather fetch",
			cityName: "London",
			country:  "UK",
			apiKey:   "test-api-key",
			mockResponse: Response{
				Name: "London",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "UK",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     15.5,
					Pressure: 1013,
					Humidity: 65,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{
						Main:        "Clouds",
						Description: "scattered clouds",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 5.2,
					Deg:   180,
				},
			},
			mockStatusCode: http.StatusOK,
			expectedResult: &schemata.FetchWeatherResponse{
				LocationName: "London",
				Country:      "UK",
				Temperature:  15.5,
				Description:  "scattered clouds",
				Humidity:     65,
				WindSpeed:    5.2,
			},
			expectedError: nil,
		},
		{
			name:           "city not found",
			cityName:       "NonExistentCity",
			country:        "XX",
			apiKey:         "test-api-key",
			mockResponse:   Response{},
			mockStatusCode: http.StatusNotFound,
			expectedResult: nil,
			expectedError:  NotFoundErr,
		},
		{
			name:           "server error",
			cityName:       "London",
			country:        "UK",
			apiKey:         "test-api-key",
			mockResponse:   Response{},
			mockStatusCode: http.StatusInternalServerError,
			expectedResult: nil,
			expectedError:  UnhandledError,
		},
		{
			name:     "empty country parameter",
			cityName: "Paris",
			country:  "",
			apiKey:   "test-api-key",
			mockResponse: Response{
				Name: "Paris",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "FR",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     20.0,
					Pressure: 1015,
					Humidity: 70,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{
						Main:        "Clear",
						Description: "clear sky",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 3.1,
					Deg:   90,
				},
			},
			mockStatusCode: http.StatusOK,
			expectedResult: &schemata.FetchWeatherResponse{
				LocationName: "Paris",
				Country:      "FR",
				Temperature:  20.0,
				Description:  "clear sky",
				Humidity:     70,
				WindSpeed:    3.1,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request parameters
				assert.Equal(t, "GET", r.Method)
				assert.Contains(t, r.URL.Query().Get("q"), tt.cityName)
				if tt.country != "" {
					assert.Contains(t, r.URL.Query().Get("q"), tt.country)
				}
				assert.Equal(t, tt.apiKey, r.URL.Query().Get("appid"))
				assert.Equal(t, "metric", r.URL.Query().Get("units"))

				// Set response status
				w.WriteHeader(tt.mockStatusCode)

				// Return mock response if status is OK
				if tt.mockStatusCode == http.StatusOK {
					responseJSON, err := json.Marshal(tt.mockResponse)
					require.NoError(t, err)
					w.Header().Set("Content-Type", "application/json")
					w.Write(responseJSON)
				}
			}))
			defer server.Close()

			// Set the base URL to our test server
			originalBaseURL := GetBaseURL()
			SetBaseURL(server.URL)
			defer SetBaseURL(originalBaseURL)

			// Create config
			config := conf.Config{}
			config.OpenWeather.ApiKey = tt.apiKey

			// Call the function
			result, err := FetchWeatherByLocation(context.Background(), tt.cityName, tt.country, config)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.LocationName, result.LocationName)
				assert.Equal(t, tt.expectedResult.Country, result.Country)
				assert.Equal(t, tt.expectedResult.Temperature, result.Temperature)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.Humidity, result.Humidity)
				assert.Equal(t, tt.expectedResult.WindSpeed, result.WindSpeed)
			}
		})
	}
}

func TestFetchWeatherByLocation_InvalidJSON(t *testing.T) {
	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"invalid": json`)) // Invalid JSON
	}))
	defer server.Close()

	// Set the base URL to our test server
	originalBaseURL := GetBaseURL()
	SetBaseURL(server.URL)
	defer SetBaseURL(originalBaseURL)

	// Create config
	config := conf.Config{}
	config.OpenWeather.ApiKey = "test-api-key"

	// Call the function
	result, err := FetchWeatherByLocation(context.Background(), "London", "UK", config)

	// Should return error due to invalid JSON
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchWeatherByLocation_NetworkError(t *testing.T) {
	// Set the base URL to an invalid server
	originalBaseURL := GetBaseURL()
	SetBaseURL("http://invalid-server-that-does-not-exist:12345")
	defer SetBaseURL(originalBaseURL)

	// Create config
	config := conf.Config{}
	config.OpenWeather.ApiKey = "test-api-key"

	// Call the function
	result, err := FetchWeatherByLocation(context.Background(), "London", "UK", config)

	// Should return error due to network failure
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSetBaseURL(t *testing.T) {
	originalURL := GetBaseURL()

	// Test setting a new URL
	newURL := "http://test-server.com"
	SetBaseURL(newURL)
	assert.Equal(t, newURL, GetBaseURL())

	// Restore original URL
	SetBaseURL(originalURL)
	assert.Equal(t, originalURL, GetBaseURL())
}
