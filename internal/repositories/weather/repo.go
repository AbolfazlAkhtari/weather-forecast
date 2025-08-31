package weather

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/schemata"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
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

func (r Repository) Update(ctx context.Context, id uuid.UUID, input map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Weather{}).Where("id = ?", id).Updates(input).Error
}

func (r Repository) PaginatedList(ctx context.Context, page int) (weathers []models.Weather, totalPage, count int64, err error) {
	offset := (page - 1) * schemata.PaginationLimit

	query := r.db.WithContext(ctx).Model(models.Weather{})

	query.Count(&count)

	result := query.Order("created_at desc").Offset(offset).Limit(schemata.PaginationLimit).Find(&weathers)
	totalPage = int64(math.Ceil(float64(count) / float64(schemata.PaginationLimit)))

	return weathers, totalPage, count, result.Error
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
