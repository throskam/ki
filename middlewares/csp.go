package middlewares

import (
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// ContentSecurityPolicy returns a middleware that adds a Content-Security-Policy header.
func ContentSecurityPolicy(policy url.Values) func(http.Handler) http.Handler {
	directives := []string{}

	for _, directive := range slices.Sorted(maps.Keys(policy)) {
		sources := policy[directive]
		directives = append(directives, directive+" "+strings.Join(sources, " "))
	}

	header := strings.Join(directives, "; ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", header)

			next.ServeHTTP(w, r)
		})
	}
}
