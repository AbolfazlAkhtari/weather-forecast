package conf

import (
	"github.com/caarlos0/env/v11"
	"log"
)

type Config struct {
	OpenWeather struct {
		ApiKey string `env:"OPEN_WEATHER_API_KEY"`
	}
}

func LoadFromEnv() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing weather api config:", err)
	}

	return cfg
}
