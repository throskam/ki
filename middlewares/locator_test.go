package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/throskam/ki"
)

func TestLocatorMiddleware(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURL := "/foo"
		actualURL := ki.GetLocation(r.Context(), "foo").URL().String()
		if actualURL != expectedURL {
			t.Errorf("expected URL %q, got %q", expectedURL, actualURL)
		}

		w.WriteHeader(http.StatusOK)
	})

	router := ki.NewMux()

	router.Get("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, ki.WithName("foo"))

	middleware := Locator(router)
	handler := middleware(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

