package weather

import (
	"encoding/json"
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

func mapUpdateInputToRepoInput(input UpdateInput) (map[string]interface{}, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var repoInput map[string]interface{}
	if err := json.Unmarshal(data, &repoInput); err != nil {
		return nil, err
	}

	return repoInput, nil
}
