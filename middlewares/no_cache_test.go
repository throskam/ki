package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNoCacheMiddleware(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := NoCache()
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedHeaders := map[string]string{
		"Cache-Control":     "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		"Pragma":            "no-cache",
		"Surrogate-Control": "no-store",
		"Expires":           time.Unix(0, 0).UTC().Format(http.TimeFormat),
	}

	for key, expected := range expectedHeaders {
		got := rec.Header().Get(key)
		if got != expected {
			t.Errorf("header %q: expected %q, got %q", key, expected, got)
		}
	}
}