package schemata

type FetchWeatherResponse struct {
	LocationName string
	Country      string
	Temperature  float64
	Description  string
	Humidity     int
	WindSpeed    float64
}
