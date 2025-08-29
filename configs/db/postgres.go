package db

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Postgres struct{}

type PgConfig struct {
	HOST     string `env:"DB_HOST"`
	PORT     string `env:"DB_PORT"`
	DATABASE string `env:"DB_DATABASE"`
	USER     string `env:"DB_USER"`
	PASSWORD string `env:"DB_PASSWORD"`
	SSLMODE  string `env:"DB_SSLMODE" envDefault:"disable"`
	TZ       string `env:"DB_TZ" envDefault:"UTC"`
}

func (Postgres) dsn() string {
	cfg := PgConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing postgres config:", err)
	}

	// postgres://user:password@host:port/dbname?sslmode=disable&TimeZone=UTC
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v&TimeZone=%v",
		cfg.USER,
		cfg.PASSWORD,
		cfg.HOST,
		cfg.PORT,
		cfg.DATABASE,
		cfg.SSLMODE,
		cfg.TZ,
	)
}

func (db Postgres) Open() gorm.Dialector {
	return postgres.Open(db.dsn())
}
