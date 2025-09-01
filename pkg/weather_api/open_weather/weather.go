package open_weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/conf"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
)

var baseURL = "https://api.openweathermap.org/data/2.5/weather"

var (
	NotFoundErr    = errors.New("not-found")
	UnhandledError = errors.New("unhandled-error")
)

// SetBaseURL allows setting the base URL for testing purposes
func SetBaseURL(url string) {
	baseURL = url
}

// GetBaseURL returns the current base URL
func GetBaseURL() string {
	return baseURL
}

func FetchWeatherByLocation(ctx context.Context, cityName, country string, config conf.Config) (*schemata.FetchWeatherResponse, error) {
	if country != "" {
		country = "," + country
	}

	url := fmt.Sprintf("%s?q=%s%s&appid=%s&units=metric", GetBaseURL(), cityName, country, config.OpenWeather.ApiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error in closing body in open_weather.FetchWeatherByLocation")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = NotFoundErr
		default:
			log.Printf("error in open_weather.FetchWeatherByLocation %v : %v - %v", resp.StatusCode, cityName, country)
			err = UnhandledError
		}

		return nil, err
	}

	var owResp Response
	if err := json.NewDecoder(resp.Body).Decode(&owResp); err != nil {
		return nil, err
	}

	dto := mapOpenWeatherResponseToFetchWeatherResponse(owResp)

	return &dto, nil
}
