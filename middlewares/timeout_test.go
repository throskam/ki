package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutMiddleware_DeadlineExceeded(t *testing.T) {
	timeout := 10 * time.Millisecond

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		// never reached.
		case <-time.After(50 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		case <-r.Context().Done():
			return
		}
	})

	middleware := Timeout(timeout)
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status %d, got %d", http.StatusGatewayTimeout, rec.Code)
	}
}

func TestTimeoutMiddleware_NoTimeout(t *testing.T) {
	timeout := 50 * time.Millisecond

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Timeout(timeout)
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
