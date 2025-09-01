package weather_api

import (
	"testing"
)

func TestLoadFetchWeatherByLocationFunc(t *testing.T) {
	tests := []struct {
		name           string
		provider       WeatherProvider
		expectedError  bool
		expectedResult bool
	}{
		{
			name:           "should return function for valid OpenWeather provider",
			provider:       OpenWeather,
			expectedError:  false,
			expectedResult: true,
		},
		{
			name:           "should return error for invalid provider",
			provider:       WeatherProvider("InvalidProvider"),
			expectedError:  true,
			expectedResult: false,
		},
		{
			name:           "should return error for empty provider",
			provider:       WeatherProvider(""),
			expectedError:  true,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadFetchWeatherByLocationFunc(tt.provider)

			if tt.expectedError {
				if err == nil {
					t.Errorf("LoadFetchWeatherByLocationFunc() expected error but got none")
				}
				if result != nil {
					t.Errorf("LoadFetchWeatherByLocationFunc() expected nil result but got %v", result)
				}
				// Check if the error message is correct
				expectedErrMsg := "invalid weather provider name | provider not implemented"
				if err.Error() != expectedErrMsg {
					t.Errorf("LoadFetchWeatherByLocationFunc() error message = %v, want %v", err.Error(), expectedErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("LoadFetchWeatherByLocationFunc() unexpected error = %v", err)
				}
				if result == nil {
					t.Errorf("LoadFetchWeatherByLocationFunc() expected function but got nil")
				}
			}
		})
	}
}

func TestFetchWeatherByLocationFuncMap(t *testing.T) {
	t.Run("should contain OpenWeather provider", func(t *testing.T) {
		if _, exists := FetchWeatherByLocationFunc[OpenWeather]; !exists {
			t.Errorf("FetchWeatherByLocationFunc map should contain OpenWeather provider")
		}
	})

	t.Run("should not contain invalid providers", func(t *testing.T) {
		invalidProviders := []WeatherProvider{
			"WeatherAPI",
			"AccuWeather",
			"",
		}

		for _, provider := range invalidProviders {
			if _, exists := FetchWeatherByLocationFunc[provider]; exists {
				t.Errorf("FetchWeatherByLocationFunc map should not contain provider: %s", provider)
			}
		}
	})
}

func TestWeatherProviderConstants(t *testing.T) {
	t.Run("should have correct OpenWeather constant value", func(t *testing.T) {
		expected := WeatherProvider("OpenWeather")
		if OpenWeather != expected {
			t.Errorf("OpenWeather constant = %v, want %v", OpenWeather, expected)
		}
	})
}

// Integration test to verify the function is returned correctly
func TestLoadFetchWeatherByLocationFuncIntegration(t *testing.T) {
	t.Run("should return function for OpenWeather", func(t *testing.T) {
		fetchFunc, err := LoadFetchWeatherByLocationFunc(OpenWeather)
		if err != nil {
			t.Fatalf("LoadFetchWeatherByLocationFunc() failed: %v", err)
		}

		if fetchFunc == nil {
			t.Fatal("LoadFetchWeatherByLocationFunc() returned nil function")
		}

		// Assert that fetchFunc is a function (not nil)
		if fetchFunc == nil {
			t.Error("Expected non-nil function, got nil")
		}
	})
}

func TestLoadFetchWeatherByLocationFuncErrorHandling(t *testing.T) {
	t.Run("should return specific error for invalid provider", func(t *testing.T) {
		invalidProvider := WeatherProvider("NonExistentProvider")

		_, err := LoadFetchWeatherByLocationFunc(invalidProvider)
		if err == nil {
			t.Fatal("Expected error for invalid provider, got nil")
		}

		// Check if it's the expected error type
		if err.Error() != "invalid weather provider name | provider not implemented" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}
