package url

import (
	"fmt"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpres"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func GetStringFromParam(r *http.Request, w http.ResponseWriter, keyName string) *string {
	param := chi.URLParam(r, keyName)
	if param == "" {
		msg := fmt.Sprintf("expected parameter %v is missing", keyName)
		httpres.SendResponse(w, http.StatusBadRequest, nil, &msg)
		return nil
	}

	return &param
}

func GetUUIDFromParam(r *http.Request, w http.ResponseWriter, keyName string) *uuid.UUID {
	param := chi.URLParam(r, keyName)

	parse, err := uuid.Parse(param)
	if err != nil {
		msg := err.Error()
		httpres.SendResponse(w, http.StatusBadRequest, nil, &msg)
		return nil
	}

	return &parse
}

// GetFromQuery todo: this error here is invalid
func GetFromQuery(r *http.Request, w http.ResponseWriter, keyName string) *string {
	token := r.URL.Query().Get(keyName)
	if token == "" {
		msg := "invalid-token"
		httpres.SendResponse(w, http.StatusBadRequest, nil, &msg)
		return nil
	}

	return &token
}
