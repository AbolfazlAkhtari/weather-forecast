package weather

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/repositories/weather"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api"
	weatherApiConf "github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/conf"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository weather.Repository
}

func NewService(db *gorm.DB) Service {
	return Service{
		db:         db,
		repository: weather.NewRepository(db),
	}
}

func (s Service) fetchData(ctx context.Context, input FetchDataInput) (*models.Weather, error) {
	fetchWeatherFunc, err := weather_api.LoadFetchWeatherByLocationFunc(weather_api.OpenWeather)
	if err != nil {
		return nil, err
	}

	conf := weatherApiConf.LoadFromEnv()
	fetchWeatherResponse, err := fetchWeatherFunc(ctx, input.CityName, input.Country, conf)
	if err != nil {
		return nil, err
	}

	weatherModel := mapFetchWeatherResponseToWeatherModel(*fetchWeatherResponse)

	err = s.repository.CreateWeather(ctx, &weatherModel)
	if err != nil {
		return nil, err
	}

	return &weatherModel, nil
}
