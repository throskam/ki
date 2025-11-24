package middlewares

import (
	"mime"
	"net/http"
	"slices"
	"strings"
)

// ContentCharset returns a middleware that checks the Content-Type header for a charset and ensures it is in the allowedCharsets.
func ContentCharset(allowedCharsets ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			charset := getCharset(r.Header.Get("Content-Type"))

			if charset == "" {
				charset = "utf-8"
			}

			if slices.Contains(allowedCharsets, charset) {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusUnsupportedMediaType)
		})
	}
}

// getCharset returns the charset from the Content-Type header.
func getCharset(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}

	for k, v := range params {
		if strings.ToLower(k) == "charset" {
			return v
		}
	}

	return ""
}
