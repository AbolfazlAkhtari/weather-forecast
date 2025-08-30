package db

import (
	"github.com/caarlos0/env/v11"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Postgres struct{}

type PgConfig struct {
	URL string `env:"DB_URL"`
}

func (Postgres) dsn() string {
	cfg := PgConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing postgres config:", err)
	}

	return cfg.URL
}

func (db Postgres) Open(migrate bool) (*gorm.DB, error) {
	dsn := db.dsn()

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if migrate {
		doMigrate(gormDB)
	}

	return gormDB, nil
}

func doMigrate(gormDB *gorm.DB) {
	if sqlDB, err := gormDB.DB(); err != nil {
		log.Println("could not get db from gormDB in db.Open -> migrate")
	} else {
		if err := goose.Up(sqlDB, "./migrations"); err != nil {
			log.Println("could not close database in db.Migrate")
		}
	}
}
