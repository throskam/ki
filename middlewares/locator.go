package middlewares

import (
	"net/http"

	"github.com/throskam/ki"
)

// Locator returns a middleware that sets the registry for the request.
func Locator(router ki.Router) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ki.SetRegistry(r.Context(), router.Registry())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
