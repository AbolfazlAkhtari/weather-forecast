package weather

import (
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/schemata"
)

type FetchDataInput struct {
	CityName string `json:"city_name" validate:"required"`
	Country  string `json:"country"`
}

type ListInput struct {
	Page int `json:"page"`
}

type ListOutput struct {
	Weathers   []models.Weather    `json:"data"`
	Pagination schemata.Pagination `json:"pagination"`
}
