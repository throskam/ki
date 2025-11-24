package middlewares

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/throskam/ki"
)

var RequestIDHeader = "X-Request-Id"

// RequestID returns a middleware that sets the request ID for the request.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx := ki.SetRequestID(r.Context(), requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
