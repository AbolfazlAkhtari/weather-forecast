package weather

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/open_weather"

	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	weatherApiSchemata "github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Weather{})
	require.NoError(t, err)

	return db
}

func setupWeatherAPIServer(t *testing.T, mockResponse *weatherApiSchemata.FetchWeatherResponse, statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "metric", r.URL.Query().Get("units"))

		// Set response status
		w.WriteHeader(statusCode)

		// Return mock response if status is OK
		if statusCode == http.StatusOK && mockResponse != nil {
			// Convert our internal format to OpenWeather API format
			owResponse := map[string]interface{}{
				"name": mockResponse.LocationName,
				"sys": map[string]interface{}{
					"country": mockResponse.Country,
				},
				"main": map[string]interface{}{
					"temp":     mockResponse.Temperature,
					"pressure": 1013,
					"humidity": mockResponse.Humidity,
				},
				"weather": []map[string]interface{}{
					{
						"main":        "Clear",
						"description": mockResponse.Description,
					},
				},
				"wind": map[string]interface{}{
					"speed": mockResponse.WindSpeed,
					"deg":   180,
				},
			}

			responseJSON, err := json.Marshal(owResponse)
			require.NoError(t, err)
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
		}
	}))

	return server
}

func TestService_paginatedList(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create test data
	testWeathers := []models.Weather{
		{CityName: "London", Country: "UK", Temperature: 20.5, Description: "Sunny", Humidity: 65, WindSpeed: 10.2, FetchedAt: time.Now()},
		{CityName: "Paris", Country: "France", Temperature: 18.0, Description: "Cloudy", Humidity: 70, WindSpeed: 8.5, FetchedAt: time.Now()},
		{CityName: "Berlin", Country: "Germany", Temperature: 15.5, Description: "Rainy", Humidity: 80, WindSpeed: 12.0, FetchedAt: time.Now()},
	}

	// Insert test data
	for _, w := range testWeathers {
		err := db.Create(&w).Error
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		page          int
		expectedCount int
		expectedError error
	}{
		{
			name:          "success with page 1",
			page:          1,
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "success with page 0 (defaults to 1)",
			page:          1,
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "success with page 2 (empty)",
			page:          2,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.paginatedList(context.Background(), tt.page)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Weathers, tt.expectedCount)
				assert.Equal(t, tt.page, result.Pagination.CurrentPage)
				assert.Equal(t, int64(3), result.Pagination.TotalCount)
			}
		})
	}
}

func TestService_latestByCityName(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create test data with multiple entries for the same city
	now := time.Now()
	oldWeather := models.Weather{
		CityName: "London", Country: "UK", Temperature: 18.0,
		Description: "Old", Humidity: 60, WindSpeed: 5.0,
		FetchedAt: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour),
	}
	newWeather := models.Weather{
		CityName: "London", Country: "UK", Temperature: 22.0,
		Description: "New", Humidity: 70, WindSpeed: 8.0,
		FetchedAt: now, CreatedAt: now,
	}

	err := db.Create(&oldWeather).Error
	require.NoError(t, err)
	err = db.Create(&newWeather).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		cityName       string
		expectedResult *models.Weather
		expectedError  error
	}{
		{
			name:           "success - returns latest weather",
			cityName:       "London",
			expectedResult: &newWeather,
			expectedError:  nil,
		},
		{
			name:           "case insensitive search",
			cityName:       "london",
			expectedResult: &newWeather,
			expectedError:  nil,
		},
		{
			name:           "city not found",
			cityName:       "NonExistent",
			expectedResult: nil,
			expectedError:  gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.latestByCityName(context.Background(), tt.cityName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.CityName, result.CityName)
				assert.Equal(t, tt.expectedResult.Temperature, result.Temperature)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
			}
		})
	}
}

func TestService_findById(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create test data
	testWeather := models.Weather{
		CityName: "London", Country: "UK", Temperature: 20.5,
		Description: "Sunny", Humidity: 65, WindSpeed: 10.2,
		FetchedAt: time.Now(),
	}

	err := db.Create(&testWeather).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		id             uuid.UUID
		expectedResult *models.Weather
		expectedError  error
	}{
		{
			name:           "success",
			id:             testWeather.ID,
			expectedResult: &testWeather,
			expectedError:  nil,
		},
		{
			name:           "id not found",
			id:             uuid.New(),
			expectedResult: nil,
			expectedError:  gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.findById(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.CityName, result.CityName)
			}
		})
	}
}

func TestService_deleteById(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create test data
	testWeather := models.Weather{
		CityName: "London", Country: "UK", Temperature: 20.5,
		Description: "Sunny", Humidity: 65, WindSpeed: 10.2,
		FetchedAt: time.Now(),
	}

	err := db.Create(&testWeather).Error
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            uuid.UUID
		expectedError error
	}{
		{
			name:          "success",
			id:            testWeather.ID,
			expectedError: nil,
		},
		{
			name:          "id not found",
			id:            uuid.New(),
			expectedError: nil, // GORM doesn't return error for non-existent delete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.deleteById(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify deletion
			if tt.name == "success" {
				var count int64
				db.Model(&models.Weather{}).Where("id = ?", tt.id).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}

func TestService_update(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create test data
	testWeather := models.Weather{
		CityName: "London", Country: "UK", Temperature: 20.5,
		Description: "Sunny", Humidity: 65, WindSpeed: 10.2,
		FetchedAt: time.Now(),
	}

	err := db.Create(&testWeather).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		id             uuid.UUID
		input          UpdateInput
		expectedResult *models.Weather
		expectedError  error
	}{
		{
			name: "success - update city name and temperature",
			id:   testWeather.ID,
			input: UpdateInput{
				CityName:    stringPtr("Updated London"),
				Temperature: float64Ptr(25.0),
				UpdatedAt:   timePtr(time.Now()),
			},
			expectedResult: &models.Weather{
				ID:          testWeather.ID,
				CityName:    "Updated London",
				Country:     "UK",
				Temperature: 25.0,
				Description: "Sunny",
				Humidity:    65,
				WindSpeed:   10.2,
			},
			expectedError: nil,
		},
		{
			name: "id not found",
			id:   uuid.New(),
			input: UpdateInput{
				CityName: stringPtr("Updated City"),
			},
			expectedResult: nil,
			expectedError:  gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.update(context.Background(), tt.id, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.CityName, result.CityName)
				assert.Equal(t, tt.expectedResult.Temperature, result.Temperature)
				assert.Equal(t, tt.expectedResult.Country, result.Country)
			}
		})
	}
}

func TestService_fetchData(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	mockResponse := &weatherApiSchemata.FetchWeatherResponse{
		LocationName: "London",
		Country:      "UK",
		Temperature:  20.5,
		Description:  "Sunny",
		Humidity:     65,
		WindSpeed:    10.2,
	}

	tests := []struct {
		name           string
		input          FetchDataInput
		mockResponse   *weatherApiSchemata.FetchWeatherResponse
		statusCode     int
		expectedResult *models.Weather
		expectedError  error
	}{
		{
			name: "success",
			input: FetchDataInput{
				CityName: "London",
				Country:  "UK",
			},
			mockResponse: mockResponse,
			statusCode:   http.StatusOK,
			expectedResult: &models.Weather{
				CityName:    "London",
				Country:     "UK",
				Temperature: 20.5,
				Description: "Sunny",
				Humidity:    65,
				WindSpeed:   10.2,
			},
			expectedError: nil,
		},
		{
			name: "API error - 404",
			input: FetchDataInput{
				CityName: "InvalidCity",
				Country:  "InvalidCountry",
			},
			mockResponse:   nil,
			statusCode:     http.StatusNotFound,
			expectedResult: nil,
			expectedError:  errors.New("not-found"),
		},
		{
			name: "API error - 500",
			input: FetchDataInput{
				CityName: "ServerError",
				Country:  "Error",
			},
			mockResponse:   nil,
			statusCode:     http.StatusInternalServerError,
			expectedResult: nil,
			expectedError:  errors.New("unhandled-error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := setupWeatherAPIServer(t, tt.mockResponse, tt.statusCode)

			// Set the base URL to our test server
			originalBaseURL := open_weather.GetBaseURL()
			open_weather.SetBaseURL(server.URL)
			defer open_weather.SetBaseURL(originalBaseURL)

			result, err := service.fetchData(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.CityName, result.CityName)
				assert.Equal(t, tt.expectedResult.Country, result.Country)
				assert.Equal(t, tt.expectedResult.Temperature, result.Temperature)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.Humidity, result.Humidity)
				assert.Equal(t, tt.expectedResult.WindSpeed, result.WindSpeed)
				assert.NotZero(t, result.FetchedAt)
				assert.NotZero(t, result.ID)
			}
		})
	}
}
