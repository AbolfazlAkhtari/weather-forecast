package httpres

import (
	"encoding/json"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"net/http"
)

type Template struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendResponse(w http.ResponseWriter, statusCode int, data any, customMessage *string) {
	var message string
	if customMessage != nil {
		message = *customMessage
	} else {
		message = http.StatusText(statusCode)
	}

	template := Template{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	w.WriteHeader(template.Code)

	response, marshalError := json.Marshal(template)
	if marshalError != nil {
		exception.ReportException(marshalError)
	} else {
		_, writeError := w.Write(response)
		if writeError != nil {
			exception.ReportException(writeError)
		}
	}
}
