package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStripePrefix(t *testing.T) {
	prefix := "/api"

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			t.Errorf("Expected path to be '/test', got '%s'", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := StripePrefix(prefix)
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}