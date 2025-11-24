package middlewares

import (
	"net/http"
	"slices"
)

// ContentEncoding returns a middleware that checks the Content-Type header for an encoding and ensures it is in the allowedContentEncodings.
func ContentEncoding(allowedContentEncodings ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			contentEncodings := r.Header["Content-Encoding"]

			for _, contentEncoding := range contentEncodings {
				found := slices.Contains(allowedContentEncodings, contentEncoding)

				if !found {
					w.WriteHeader(http.StatusUnsupportedMediaType)
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
