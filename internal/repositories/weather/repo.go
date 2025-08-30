package weather

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) CreateWeather(ctx context.Context, weather *models.Weather) error {
	return r.db.WithContext(ctx).Create(weather).Error
}
