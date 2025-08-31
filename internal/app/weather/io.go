package weather

import (
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/schemata"
	"time"
)

type FetchDataInput struct {
	CityName string `json:"city_name" validate:"required"`
	Country  string `json:"country"`
}

type UpdateInput struct {
	CityName    *string    `json:"city_name,omitempty" validate:"omitempty,min=1"`
	Country     *string    `json:"country,omitempty" validate:"omitempty,alpha"`
	Temperature *float64   `json:"temperature,omitempty" validate:"omitempty"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=200"`
	Humidity    *int       `json:"humidity,omitempty" validate:"omitempty,gte=0,lte=100"`
	WindSpeed   *float64   `json:"wind_speed,omitempty" validate:"omitempty,gte=0"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ListInput struct {
	Page int `json:"page"`
}

type ListOutput struct {
	Weathers   []models.Weather    `json:"data"`
	Pagination schemata.Pagination `json:"pagination"`
}
