package http

import (
	"errors"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/open_weather"
	"gorm.io/gorm"
	"net/http"
)

func MapErrorToHttpStatusCode(err error) int {
	switch {
	case errors.Is(err, open_weather.NotFoundErr), errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	case errors.Is(err, open_weather.UnhandledError):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
