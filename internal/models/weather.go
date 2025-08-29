package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Weather struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	CityName    string    `gorm:"type:varchar(255);not null;column:city_name" json:"city_name"`
	Country     string    `gorm:"type:varchar(255);not null;column:country" json:"country"`
	Temperature float64   `gorm:"not null;column:temperature" json:"temperature"`
	Description string    `gorm:"type:varchar(255);column:description" json:"description"`
	Humidity    int       `gorm:"not null;column:humidity" json:"humidity"`
	WindSpeed   float64   `gorm:"not null;column:wind_speed" json:"wind_speed"`
	FetchedAt   time.Time `gorm:"not null;column:fetched_at" json:"fetched_at"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (w *Weather) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	return
}
