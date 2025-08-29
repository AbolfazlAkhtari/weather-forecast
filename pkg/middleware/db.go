package middleware

import (
	"context"
	"github.com/AbolfazlAkhtari/weather-forecast/pkg/httpreq"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func SetDB(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ctx := httpreq.PutDB(r, db, timeoutContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
