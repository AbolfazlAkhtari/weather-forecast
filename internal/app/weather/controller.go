package weather

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	router *chi.Mux
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		db: db,
	}
}

func (c Controller) InitRoutes() {
	c.router.Group(func(router chi.Router) {
		router.Get("/weather", nil)
		router.Get("/weather/:id", nil)
		router.Get("/weather/latest/:cityName", nil)
		router.Post("/weather", nil)
		router.Put("/weather/:id", nil)
		router.Delete("/weather/:id", nil)
	})
}
