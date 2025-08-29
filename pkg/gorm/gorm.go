package gorm

import (
	"github.com/AbolfazlAkhtari/weather-forecast/configs/db"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"gorm.io/gorm"
)

func InitializeGorm(database db.DB) *gorm.DB {
	gormDB, err := gorm.Open(database.Open(), &gorm.Config{})
	if err != nil {
		exception.ReportException(err)
		return nil
	}

	return gormDB
}
