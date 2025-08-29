package httpreq

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/exception"
	"gorm.io/gorm"
	"net/http"
)

const DBContextKey = "DB"

func GetDB(r *http.Request) *gorm.DB {
	db, ok := r.Context().Value(DBContextKey).(*gorm.DB)
	if !ok {
		exception.ReportException("Could not get DB from context")
		return nil
	}

	return db
}

func PutDB(r *http.Request, db *gorm.DB, timeoutContext context.Context) context.Context {
	ctx := context.WithValue(r.Context(), DBContextKey, db.WithContext(timeoutContext))
	return ctx
}
