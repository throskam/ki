package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovererMiddleware(t *testing.T) {
	var logBuf bytes.Buffer

	loggerHandler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(loggerHandler)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	middleware := Recoverer(logger)
	handler := middleware(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	expected := "Internal Server Error\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}

	logged := logBuf.String()
	if !strings.Contains(logged, "panic recovered") {
		t.Error("expected log to contain 'panic recovered'")
	}
	if !strings.Contains(logged, "something went wrong") {
		t.Error("expected log to contain panic error message")
	}
	if !strings.Contains(logged, "goroutine") { // stack trace usually contains this
		t.Error("expected log to contain stack trace")
	}
}