package weather

import (
	"errors"
	"fmt"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpreq"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpres"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/url"
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
		router.Get("/weather/latest/{city_name}", c.getByCityName)
		router.Get("/weather/:id", nil)
		router.Post("/weather", c.fetchData)
		router.Put("/weather/:id", nil)
		router.Delete("/weather/:id", nil)
	})
}

// todo: refactor error http status mappings

func (c Controller) getByCityName(w http.ResponseWriter, r *http.Request) {
	cityName := url.GetStringFromParam(r, w, "city_name")

	output, err := c.service.latestByCityName(r.Context(), *cityName)
	if err != nil {
		fmt.Println(output, err)
		var status int

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}

		msg := err.Error()
		httpres.SendResponse(w, status, nil, &msg)
		return
	}

	httpres.SendResponse(w, http.StatusOK, output, nil)
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
