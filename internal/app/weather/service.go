package weather

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/repositories/weather"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/schemata"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api"
	weatherApiConf "github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/conf"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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

func (s Service) paginatedList(ctx context.Context, page int) (*ListOutput, error) {
	if page == 0 {
		page = 1
	}

	weathers, totalPage, count, err := s.repository.PaginatedList(ctx, page)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Weathers: weathers,
		Pagination: schemata.Pagination{
			TotalPage:   totalPage,
			TotalCount:  count,
			CurrentPage: page,
		},
	}, nil
}

func (s Service) latestByCityName(ctx context.Context, cityName string) (*models.Weather, error) {
	w, err := s.repository.LatestByCityName(ctx, cityName)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s Service) findById(ctx context.Context, id uuid.UUID) (*models.Weather, error) {
	w, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s Service) deleteById(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteById(ctx, id)
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

	w := mapFetchWeatherResponseToWeatherModel(*fetchWeatherResponse)

	err = s.repository.Create(ctx, &w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s Service) update(ctx context.Context, id uuid.UUID, input UpdateInput) (*models.Weather, error) {
	now := time.Now()
	input.UpdatedAt = &now

	repoInput, err := mapUpdateInputToRepoInput(input)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(ctx, id, repoInput)
	if err != nil {
		return nil, err
	}

	w, err := s.findById(ctx, id)
	if err != nil {
		return nil, err
	}

	return w, nil
}
