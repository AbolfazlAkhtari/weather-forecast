package httpreq

import (
	"encoding/json"
	"errors"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpres"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/validation"
	"io"
	"net/http"
)

func ParseAndValidateInput[T any](w http.ResponseWriter, r *http.Request) (input *T) {
	err := parse(w, r, &input)
	if err != nil {
		return nil
	}

	err = validate(w, input)
	if err != nil {
		return nil
	}

	return input
}

func parse(w http.ResponseWriter, r *http.Request, input any) error {
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		exception.ReportException(err)

		var message string

		if err == io.EOF {
			message = "no input provided"
		} else {
			message = err.Error()
		}

		httpres.SendResponse(w, http.StatusBadRequest, nil, &message)
		return err
	}

	return nil
}

func validate(w http.ResponseWriter, input any) (err error) {
	validationErrors := validation.ValidateData(input)
	if validationErrors != nil {
		httpres.SendResponse(w, http.StatusUnprocessableEntity, validationErrors, nil)
		return errors.New("validation_failed")
	}

	return nil
}
