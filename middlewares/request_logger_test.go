package middlewares

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/throskam/ki"
)

func TestRequestLogger(t *testing.T) {
	var logBuf bytes.Buffer

	logHandler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(logHandler)

	ki.Logger = logger

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestLogger()
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/hello?world=true", nil)
	req.RemoteAddr = "127.0.0.1:5555"

	ctx := ki.SetRequestID(context.Background(), "test-id-123")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	logOutput := logBuf.String()

	if !strings.Contains(logOutput, "request") {
		t.Errorf("expected log to contain 'request'")
	}
	if !strings.Contains(logOutput, `method=GET`) {
		t.Errorf("expected log to contain 'method=GET'")
	}
	if !strings.Contains(logOutput, `path="/hello?world=true"`) {
		t.Errorf("expected log to contain correct path")
	}
	if !strings.Contains(logOutput, `remote=127.0.0.1:5555`) {
		t.Errorf("expected log to contain correct remote")
	}
	if !strings.Contains(logOutput, `status=200`) {
		t.Errorf("expected log to contain status=200")
	}
	if !strings.Contains(logOutput, `size=`) {
		t.Errorf("expected log to contain size")
	}
	if !strings.Contains(logOutput, `duration=`) {
		t.Errorf("expected log to contain duration")
	}
	if !strings.Contains(logOutput, `requestID=test-id-123`) {
		t.Errorf("expected log to contain correct requestID")
	}
}
