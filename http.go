package isulog

import (
	"log"
	"net/http"
	"time"

	"github.com/pirosiki197/isulog/internal"
)

func StdHTTP() func(next http.Handler) http.Handler {
	return StdHTTPWithConfig(DefaultConfig)
}

func StdHTTPWithConfig(config Config) func(next http.Handler) http.Handler {
	recorder := internal.NewRecorder(config.Filename)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			record := internal.Record{
				Path:         r.Pattern,
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
