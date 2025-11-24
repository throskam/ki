package middlewares

import (
	"net"
	"net/http"
	"strings"
)

// RealIP returns a middleware that sets the remote address for the request.
func RealIP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.RemoteAddr = getRealIP(r)

			next.ServeHTTP(w, r)
		})
	}
}

// getRealIP returns the remote address for the request.
func getRealIP(r *http.Request) string {
	ips := []string{
		r.Header.Get("True-Client-IP"),
		r.Header.Get("X-Real-IP"),
		r.Header.Get("X-Forwarded-For"),
	}

	for _, ip := range ips {
		// X-Forwarded-For contains a list of IPs separated by commas.
		ip, _, _ = strings.Cut(ip, ",")

		if ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}

	return r.RemoteAddr
}
