package weather

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/google/uuid"
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

func (r Repository) Create(ctx context.Context, w *models.Weather) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r Repository) LatestByCityName(ctx context.Context, cityName string) (w *models.Weather, err error) {
	err = r.db.WithContext(ctx).Where("LOWER(city_name) = LOWER(?)", cityName).Order("created_at DESC").First(&w).Error

	return w, err
}

func (r Repository) FindById(ctx context.Context, id uuid.UUID) (w *models.Weather, err error) {
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&w).Error

	return w, err
}

func (r Repository) DeleteById(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(models.Weather{}, id).Error
}
