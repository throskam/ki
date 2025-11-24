package middlewares

import "net/http"

// StripePrefix returns a middleware that strips the prefix from the request path.
func StripePrefix(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}