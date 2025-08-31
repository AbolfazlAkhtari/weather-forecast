package weather

import (
	httpErr "github.com/AbolfazlAkhtari/weather-forecast/internal/pkg/http"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpreq"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpres"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/url"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
		router.Get("/weather", c.paginatedList)
		router.Get("/weather/latest/{city_name}", c.getByCityName)
		router.Get("/weather/{id}", c.getById)
		router.Post("/weather", c.fetchData)
		router.Put("/weather/{id}", nil)
		router.Delete("/weather/{id}", c.deleteById)
	})
}

func (c Controller) paginatedList(w http.ResponseWriter, r *http.Request) {
	page := 1

	pageInput := r.URL.Query().Get("page")
	if pageInput != "" {
		var err error
		page, err = strconv.Atoi(pageInput)

		if err != nil {
			msg := err.Error()
			httpres.SendResponse(w, http.StatusBadRequest, nil, &msg)
			return
		}
	}

	output, err := c.service.paginatedList(r.Context(), page)
	if err != nil {
		handleServiceErrors(w, err)
		return
	}

	httpres.SendResponse(w, http.StatusOK, output, nil)
}

func (c Controller) getByCityName(w http.ResponseWriter, r *http.Request) {
	cityName := url.GetStringFromParam(r, w, "city_name")
	if cityName == nil {
		return
	}

	output, err := c.service.latestByCityName(r.Context(), *cityName)
	if err != nil {
		handleServiceErrors(w, err)
		return
	}

	httpres.SendResponse(w, http.StatusOK, output, nil)
}

func (c Controller) getById(w http.ResponseWriter, r *http.Request) {
	id := url.GetUUIDFromParam(r, w, "id")
	if id == nil {
		return
	}

	output, err := c.service.findById(r.Context(), *id)
	if err != nil {
		handleServiceErrors(w, err)
		return
	}

	httpres.SendResponse(w, http.StatusOK, output, nil)
}

func (c Controller) deleteById(w http.ResponseWriter, r *http.Request) {
	id := url.GetUUIDFromParam(r, w, "id")
	if id == nil {
		return
	}

	err := c.service.deleteById(r.Context(), *id)
	if err != nil {
		handleServiceErrors(w, err)
		return
	}

	httpres.SendResponse(w, http.StatusOK, nil, nil)
}

func (c Controller) fetchData(w http.ResponseWriter, r *http.Request) {
	input := httpreq.ParseAndValidateInput[FetchDataInput](w, r)
	if input == nil {
		return
	}

	output, err := c.service.fetchData(r.Context(), *input)
	if err != nil {
		handleServiceErrors(w, err)
		return
	}

	httpres.SendResponse(w, http.StatusCreated, output, nil)
}

func handleServiceErrors(w http.ResponseWriter, err error) {
	status := httpErr.MapErrorToHttpStatusCode(err)

	msg := err.Error()
	httpres.SendResponse(w, status, nil, &msg)
}
