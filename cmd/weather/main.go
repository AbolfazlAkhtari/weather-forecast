package main

import (
	"fmt"
	"github.com/AbolfazlAkhtari/weather-forecast/configs/db"
	weatherCfg "github.com/AbolfazlAkhtari/weather-forecast/configs/weather"
	"github.com/AbolfazlAkhtari/weather-forecast/internal/app/weather"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			exception.ReportException(r)
		}
	}()

	database, err := db.Postgres{}.Open(true)
	if err != nil {
		exception.ReportException(err)
		return
	}

	config := weatherCfg.LoadFromEnv()

	httpRouter := chi.NewRouter()
	httpRouter.Use(chiMiddleware.Logger)
	httpRouter.Use(middleware.SetResponseHeader)

	httpRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.AllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
	}))

	weather.NewController(database, httpRouter).InitRoutes()

	fmt.Printf("App Served on port %v \n\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), httpRouter))
}
