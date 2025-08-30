package weather

import (
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpreq"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpres"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/weather_api/open_weather"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type Controller struct {
	db      *gorm.DB
	service Service
	router  *chi.Mux
}

func NewController(db *gorm.DB, router *chi.Mux) Controller {
	return Controller{
		db:      db,
		service: NewService(db),
		router:  router,
	}
}

func (c Controller) InitRoutes() {
	c.router.Group(func(router chi.Router) {
		router.Get("/weather", nil)
		router.Get("/weather/:id", nil)
		router.Get("/weather/latest/:city_name", nil)
		router.Post("/weather", c.fetchData)
		router.Put("/weather/:id", nil)
		router.Delete("/weather/:id", nil)
	})
}

func (c Controller) fetchData(w http.ResponseWriter, r *http.Request) {
	input := httpreq.ParseAndValidateInput[FetchDataInput](w, r)
	if input == nil {
		return
	}

	output, err := c.service.fetchData(r.Context(), *input)
	if err != nil {
		var status int

		switch err {
		case open_weather.NotFoundErr:
			status = http.StatusNotFound
		case open_weather.UnhandledError:
			status = http.StatusServiceUnavailable
		}

		msg := err.Error()
		httpres.SendResponse(w, status, nil, &msg)
		return
	}

	httpres.SendResponse(w, http.StatusCreated, output, nil)
}
