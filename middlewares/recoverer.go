package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recoverer returns a middleware that recovers from panics and returns an Internal Server Error.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack", string(stack)),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}