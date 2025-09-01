package open_weather

import (
	"testing"

	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
	"github.com/stretchr/testify/assert"
)

func TestMapOpenWeatherResponseToFetchWeatherResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    Response
		expected schemata.FetchWeatherResponse
	}{
		{
			name: "complete weather data",
			input: Response{
				Name: "London",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "GB",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     15.5,
					Pressure: 1013,
					Humidity: 75,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{
						Main:        "Clouds",
						Description: "scattered clouds",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 5.2,
					Deg:   180,
				},
			},
			expected: schemata.FetchWeatherResponse{
				LocationName: "London",
				Country:      "GB",
				Temperature:  15.5,
				Description:  "scattered clouds",
				Humidity:     75,
				WindSpeed:    5.2,
			},
		},
		{
			name: "empty weather array",
			input: Response{
				Name: "Paris",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "FR",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     22.0,
					Pressure: 1015,
					Humidity: 60,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 3.1,
					Deg:   90,
				},
			},
			expected: schemata.FetchWeatherResponse{
				LocationName: "Paris",
				Country:      "FR",
				Temperature:  22.0,
				Description:  "",
				Humidity:     60,
				WindSpeed:    3.1,
			},
		},
		{
			name: "multiple weather conditions",
			input: Response{
				Name: "New York",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "US",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     -5.0,
					Pressure: 1020,
					Humidity: 85,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{
						Main:        "Snow",
						Description: "light snow",
					},
					{
						Main:        "Clouds",
						Description: "overcast clouds",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 8.5,
					Deg:   270,
				},
			},
			expected: schemata.FetchWeatherResponse{
				LocationName: "New York",
				Country:      "US",
				Temperature:  -5.0,
				Description:  "light snow",
				Humidity:     85,
				WindSpeed:    8.5,
			},
		},
		{
			name: "zero values",
			input: Response{
				Name: "Tokyo",
				Sys: struct {
					Country string `json:"country"`
				}{
					Country: "JP",
				},
				Main: struct {
					Temp     float64 `json:"temp"`
					Pressure int     `json:"pressure"`
					Humidity int     `json:"humidity"`
				}{
					Temp:     0.0,
					Pressure: 0,
					Humidity: 0,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{
					{
						Main:        "Clear",
						Description: "clear sky",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 0.0,
					Deg:   0,
				},
			},
			expected: schemata.FetchWeatherResponse{
				LocationName: "Tokyo",
				Country:      "JP",
				Temperature:  0.0,
				Description:  "clear sky",
				Humidity:     0,
				WindSpeed:    0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapOpenWeatherResponseToFetchWeatherResponse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}