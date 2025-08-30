package weather_api

import (
	"context"
	"errors"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/conf"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/open_weather"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
)

type WeatherProvider string

const (
	OpenWeather WeatherProvider = "OpenWeather"
)

var FetchWeatherByLocationFunc = map[WeatherProvider]FetchWeatherByLocation{
	OpenWeather: open_weather.FetchWeatherByLocation,
}

type FetchWeatherByLocation func(ctx context.Context, cityName, country string, config conf.Config) (*schemata.FetchWeatherResponse, error)

func LoadFetchWeatherByLocationFunc(provider WeatherProvider) (FetchWeatherByLocation, error) {
	fetchWeather, ok := FetchWeatherByLocationFunc[provider]
	if !ok {
		return nil, errors.New("invalid weather provider name | provider not implemented")
	}

	return fetchWeather, nil
}
