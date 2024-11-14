package isulog

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pirosiki197/isulog/internal"
)

func Chi() func(next http.Handler) http.Handler {
	return ChiWithConfig(DefaultConfig)
}

func ChiWithConfig(config Config) func(next http.Handler) http.Handler {
	recorder := internal.NewRecorder(config.Filename)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			record := internal.Record{
				Path:         chi.RouteContext(r.Context()).RoutePattern(),
				Method:       r.Method,
				StatusCode:   200, // TODO: capture status code
				ResponseTime: time.Since(start),
			}
			if err := recorder.Save(record); err != nil {
				log.Printf("[ERROR] Failed to record. isulog: %s", err.Error())
			}
		}
		return http.HandlerFunc(fn)
	}
}
