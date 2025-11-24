package middlewares

import (
	"net/http"
	"time"
)

// NoCache returns a middleware that sets the appropriate headers to disable caching.
func NoCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Expires", time.Unix(0, 0).UTC().Format(http.TimeFormat))
			w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Surrogate-Control", "no-store")

			next.ServeHTTP(w, r)
		})
	}
}