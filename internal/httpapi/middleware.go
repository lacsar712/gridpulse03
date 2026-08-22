package httpapi

import (
	"log"
	"net/http"
	"time"

	"example.com/gridpulse/internal/clock"
)

func withTimeout(h http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(h, d, "timeout")
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %dms", r.Method, r.URL.Path, clock.SinceMS(start))
	})
}

func jsonOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			ct := r.Header.Get("Content-Type")
			if ct != "" && ct != "application/json" && ct != "application/json; charset=utf-8" {
				http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
