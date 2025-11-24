package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/throskam/ki"
)

// RequestLogger returns a middleware that logs the request.
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			start := time.Now()

			logger := ki.Logger.With(
				slog.String("requestID", ki.GetRequestID(ctx)),
			)

			ctx = ki.SetLogger(ctx, logger)

			brw := ki.NewBufferedResponseWriter(w)

			next.ServeHTTP(brw, r.WithContext(ctx))

			end := time.Now()
			duration := end.Sub(start)

			ki.MustGetLogger(ctx).LogAttrs(
				ctx,
				slog.LevelInfo,
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.RequestURI()),
				slog.String("remote", r.RemoteAddr),
				slog.Int("status", brw.StatusCode()),
				slog.Int("size", brw.Size()),
				slog.Int64("duration", duration.Microseconds()),
			)

			brw.Flush()
		})
	}
}

