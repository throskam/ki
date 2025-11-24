package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/throskam/ki"
)

func TestRequestIDMiddleware_HeaderPresent(t *testing.T) {
	expectedID := "test-request-id"

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID := ki.GetRequestID(r.Context())
		if ctxID != expectedID {
			t.Errorf("expected request ID %q, got %q", expectedID, ctxID)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	req.Header.Set(RequestIDHeader, expectedID)

	middleware := RequestID()
	handler := middleware(noopHandler)

	handler.ServeHTTP(rec, req)
}

func TestRequestIDMiddleware_HeaderMissing(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ki.GetRequestID(r.Context())
		if requestID == "" {
			t.Error("expected a generated request ID, got empty string")
		}

		_, err := uuid.Parse(requestID)
		if err != nil {
			t.Errorf("expected a valid UUID, got %q", requestID)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware := RequestID()
	handler := middleware(noopHandler)

	handler.ServeHTTP(rec, req)
}

