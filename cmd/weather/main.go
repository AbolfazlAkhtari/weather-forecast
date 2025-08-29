package main

import (
	"fmt"
	"github.com/AbolfazlAkhtari/weather-forecast/configs/db"
	"github.com/AbolfazlAkhtari/weather-forecast/configs/weather"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/gorm"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			exception.ReportException(r)
		}
	}()

	database := gorm.InitializeGorm(db.Postgres{})

	httpRouter := chi.NewRouter()
	httpRouter.Use(chiMiddleware.Logger)
	httpRouter.Use(middleware.SetResponseHeader)
	httpRouter.Use(func(handler http.Handler) http.Handler {
		return middleware.SetDB(handler, database)
	})

	//auth.Controller{ChiRouter: httpRouter}.InitRoutes()

	config := weather.LoadFromEnv()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.PORT), httpRouter))
}
