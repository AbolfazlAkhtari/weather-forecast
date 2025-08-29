package weather

import (
	"github.com/caarlos0/env/v11"
	"log"
)

type Config struct {
	PORT string `env:"WEATHER_PORT"`
}

func LoadFromEnv() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing weather config:", err)
	}

	return cfg
}
