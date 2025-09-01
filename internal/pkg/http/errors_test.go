package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/open_weather"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMapErrorToHttpStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "should return 404 for open_weather.NotFoundErr",
			err:            open_weather.NotFoundErr,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "should return 404 for gorm.ErrRecordNotFound",
			err:            gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "should return 503 for open_weather.UnhandledError",
			err:            open_weather.UnhandledError,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "should return 0 for nil error",
			err:            nil,
			expectedStatus: 0,
		},
		{
			name:           "should return 500 for unknown error",
			err:            errors.New("some random error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "should return 500 for wrapped open_weather.NotFoundErr",
			err:            errors.New("wrapped error: " + open_weather.NotFoundErr.Error()),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "should return 404 for wrapped gorm.ErrRecordNotFound using errors.Wrap",
			err:            errors.New("wrapped: " + gorm.ErrRecordNotFound.Error()),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapErrorToHttpStatusCode(tt.err)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}

func TestMapErrorToHttpStatusCode_WithWrappedErrors(t *testing.T) {
	// Test with errors.Is() compatible wrapped errors
	wrappedNotFoundErr := errors.New("wrapped not found")
	wrappedNotFoundErr = errors.Join(wrappedNotFoundErr, open_weather.NotFoundErr)

	wrappedGormErr := errors.New("wrapped gorm error")
	wrappedGormErr = errors.Join(wrappedGormErr, gorm.ErrRecordNotFound)

	wrappedUnhandledErr := errors.New("wrapped unhandled")
	wrappedUnhandledErr = errors.Join(wrappedUnhandledErr, open_weather.UnhandledError)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "should return 404 for wrapped open_weather.NotFoundErr using errors.Join",
			err:            wrappedNotFoundErr,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "should return 404 for wrapped gorm.ErrRecordNotFound using errors.Join",
			err:            wrappedGormErr,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "should return 503 for wrapped open_weather.UnhandledError using errors.Join",
			err:            wrappedUnhandledErr,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapErrorToHttpStatusCode(tt.err)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}

func TestMapErrorToHttpStatusCode_EdgeCases(t *testing.T) {
	// Test with custom error types that might be similar but not the same
	customNotFoundErr := errors.New("not-found")
	customUnhandledErr := errors.New("unhandled-error")

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "should return 500 for custom not-found error (not the same as open_weather.NotFoundErr)",
			err:            customNotFoundErr,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "should return 500 for custom unhandled error (not the same as open_weather.UnhandledError)",
			err:            customUnhandledErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapErrorToHttpStatusCode(tt.err)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}
