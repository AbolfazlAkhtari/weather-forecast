package open_weather

import (
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/schemata"
)

func mapOpenWeatherResponseToFetchWeatherResponse(owResp Response) schemata.FetchWeatherResponse {
	resp := schemata.FetchWeatherResponse{
		LocationName: owResp.Name,
		Country:      owResp.Sys.Country,
		Temperature:  owResp.Main.Temp,
		Description:  "",
		Humidity:     owResp.Main.Humidity,
		WindSpeed:    owResp.Wind.Speed,
	}

	if len(owResp.Weather) > 0 {
		resp.Description = owResp.Weather[0].Description
	}

	return resp
}
