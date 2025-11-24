package middlewares

import (
	"net/http"
	"slices"
	"strings"
)

// ContentType returns a middleware that checks the Content-Type header for a type and ensures it is in the allowedContentTypes.
func ContentType(allowedContentTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			contentType, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";")

			if slices.Contains(allowedContentTypes, contentType) {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusUnsupportedMediaType)
		})
	}
}