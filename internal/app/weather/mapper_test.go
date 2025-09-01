package weather

import (
	"testing"
	"time"

	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
	"github.com/stretchr/testify/assert"
)

func TestMapFetchWeatherResponseToWeatherModel(t *testing.T) {
	tests := []struct {
		name     string
		response schemata.FetchWeatherResponse
		expected models.Weather
	}{
		{
			name: "successful mapping with all fields",
			response: schemata.FetchWeatherResponse{
				LocationName: "New York",
				Country:      "US",
				Temperature:  25.5,
				Description:  "Sunny",
				Humidity:     65,
				WindSpeed:    10.2,
			},
			expected: models.Weather{
				CityName:    "New York",
				Country:     "US",
				Temperature: 25.5,
				Description: "Sunny",
				Humidity:    65,
				WindSpeed:   10.2,
			},
		},
		{
			name: "successful mapping with zero values",
			response: schemata.FetchWeatherResponse{
				LocationName: "",
				Country:      "",
				Temperature:  0.0,
				Description:  "",
				Humidity:     0,
				WindSpeed:    0.0,
			},
			expected: models.Weather{
				CityName:    "",
				Country:     "",
				Temperature: 0.0,
				Description: "",
				Humidity:    0,
				WindSpeed:   0.0,
			},
		},
		{
			name: "successful mapping with negative values",
			response: schemata.FetchWeatherResponse{
				LocationName: "Moscow",
				Country:      "RU",
				Temperature:  -15.3,
				Description:  "Snowy",
				Humidity:     80,
				WindSpeed:    5.7,
			},
			expected: models.Weather{
				CityName:    "Moscow",
				Country:     "RU",
				Temperature: -15.3,
				Description: "Snowy",
				Humidity:    80,
				WindSpeed:   5.7,
			},
		},
		{
			name: "successful mapping with special characters",
			response: schemata.FetchWeatherResponse{
				LocationName: "São Paulo",
				Country:      "BR",
				Temperature:  30.1,
				Description:  "Partly Cloudy",
				Humidity:     70,
				WindSpeed:    8.9,
			},
			expected: models.Weather{
				CityName:    "São Paulo",
				Country:     "BR",
				Temperature: 30.1,
				Description: "Partly Cloudy",
				Humidity:    70,
				WindSpeed:   8.9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapFetchWeatherResponseToWeatherModel(tt.response)

			// Check all fields except FetchedAt since it's set to time.Now()
			assert.Equal(t, tt.expected.CityName, result.CityName)
			assert.Equal(t, tt.expected.Country, result.Country)
			assert.Equal(t, tt.expected.Temperature, result.Temperature)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Humidity, result.Humidity)
			assert.Equal(t, tt.expected.WindSpeed, result.WindSpeed)

			// Check that FetchedAt is set to a recent time
			assert.WithinDuration(t, time.Now(), result.FetchedAt, 2*time.Second)
		})
	}
}

func TestMapUpdateInputToRepoInput(t *testing.T) {
	tests := []struct {
		name        string
		input       UpdateInput
		expectError bool
		expected    map[string]interface{}
	}{
		{
			name: "successful mapping with all fields",
			input: UpdateInput{
				CityName:    stringPtr("London"),
				Country:     stringPtr("UK"),
				Temperature: float64Ptr(18.5),
				Description: stringPtr("Rainy"),
				Humidity:    intPtr(85),
				WindSpeed:   float64Ptr(12.3),
				UpdatedAt:   timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			},
			expectError: false,
			expected: map[string]interface{}{
				"city_name":   "London",
				"country":     "UK",
				"temperature": 18.5,
				"description": "Rainy",
				"humidity":    float64(85), // JSON unmarshals int to float64
				"wind_speed":  12.3,
				"updated_at":  "2024-01-15T10:30:00Z",
			},
		},
		{
			name: "successful mapping with partial fields",
			input: UpdateInput{
				CityName:  stringPtr("Paris"),
				Country:   stringPtr("FR"),
				Humidity:  intPtr(60),
				UpdatedAt: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			},
			expectError: false,
			expected: map[string]interface{}{
				"city_name":  "Paris",
				"country":    "FR",
				"humidity":   float64(60),
				"updated_at": "2024-01-15T10:30:00Z",
			},
		},
		{
			name: "successful mapping with nil fields",
			input: UpdateInput{
				CityName:    nil,
				Country:     nil,
				Temperature: nil,
				Description: nil,
				Humidity:    nil,
				WindSpeed:   nil,
				UpdatedAt:   nil,
			},
			expectError: false,
			expected:    map[string]interface{}{},
		},
		{
			name: "successful mapping with zero values",
			input: UpdateInput{
				CityName:    stringPtr(""),
				Country:     stringPtr(""),
				Temperature: float64Ptr(0.0),
				Description: stringPtr(""),
				Humidity:    intPtr(0),
				WindSpeed:   float64Ptr(0.0),
				UpdatedAt:   timePtr(time.Time{}),
			},
			expectError: false,
			expected: map[string]interface{}{
				"city_name":   "",
				"country":     "",
				"temperature": 0.0,
				"description": "",
				"humidity":    float64(0),
				"wind_speed":  0.0,
				"updated_at":  "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "successful mapping with negative values",
			input: UpdateInput{
				CityName:    stringPtr("Oslo"),
				Country:     stringPtr("NO"),
				Temperature: float64Ptr(-5.2),
				Description: stringPtr("Cold"),
				Humidity:    intPtr(90),
				WindSpeed:   float64Ptr(15.7),
				UpdatedAt:   timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			},
			expectError: false,
			expected: map[string]interface{}{
				"city_name":   "Oslo",
				"country":     "NO",
				"temperature": -5.2,
				"description": "Cold",
				"humidity":    float64(90),
				"wind_speed":  15.7,
				"updated_at":  "2024-01-15T10:30:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapUpdateInputToRepoInput(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Compare the result with expected values
				for key, expectedValue := range tt.expected {
					actualValue, exists := result[key]
					assert.True(t, exists, "Key %s should exist in result", key)
					assert.Equal(t, expectedValue, actualValue, "Value for key %s should match", key)
				}

				// Check that no extra keys exist in result
				for key := range result {
					_, exists := tt.expected[key]
					assert.True(t, exists, "Key %s should not exist in result", key)
				}
			}
		})
	}
}

func TestMapUpdateInputToRepoInput_EdgeCases(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		input := UpdateInput{}
		result, err := mapUpdateInputToRepoInput(input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("mixed nil and non-nil fields", func(t *testing.T) {
		input := UpdateInput{
			CityName:    stringPtr("Tokyo"),
			Country:     nil,
			Temperature: float64Ptr(22.0),
			Description: nil,
			Humidity:    intPtr(70),
			WindSpeed:   nil,
			UpdatedAt:   timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		}
		result, err := mapUpdateInputToRepoInput(input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 4) // Only non-nil fields should be included

		expected := map[string]interface{}{
			"city_name":   "Tokyo",
			"temperature": 22.0,
			"humidity":    float64(70),
			"updated_at":  "2024-01-15T10:30:00Z",
		}

		for key, expectedValue := range expected {
			actualValue, exists := result[key]
			assert.True(t, exists, "Key %s should exist in result", key)
			assert.Equal(t, expectedValue, actualValue, "Value for key %s should match", key)
		}
	})
}

// Helper functions for creating pointers to values
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
