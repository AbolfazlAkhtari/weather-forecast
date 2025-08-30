package weather

import (
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
	"time"
)

func mapFetchWeatherResponseToWeatherModel(response schemata.FetchWeatherResponse) models.Weather {
	return models.Weather{
		CityName:    response.LocationName,
		Country:     response.Country,
		Temperature: response.Temperature,
		Description: response.Description,
		Humidity:    response.Humidity,
		WindSpeed:   response.WindSpeed,
		FetchedAt:   time.Now(),
	}
}
