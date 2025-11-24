package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentTypeMiddleware(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name                string
		contentType         string
		contentLength       int64
		allowedContentTypes []string
		expectedStatusCode  int
	}{
		{
			name:                "Valid content type - application/json",
			contentType:         "application/json",
			contentLength:       10,
			allowedContentTypes: []string{"application/json"},
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                "Invalid content type - text/plain",
			contentType:         "text/plain",
			contentLength:       10,
			allowedContentTypes: []string{"application/json"},
			expectedStatusCode:  http.StatusUnsupportedMediaType,
		},
		{
			name:                "No content length - should pass",
			contentType:         "text/plain",
			contentLength:       0,
			allowedContentTypes: []string{"application/json"},
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                "Content type with charset - should still match",
			contentType:         "application/json; charset=utf-8",
			contentLength:       10,
			allowedContentTypes: []string{"application/json"},
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                "Empty content type header - should reject if body exists",
			contentType:         "",
			contentLength:       10,
			allowedContentTypes: []string{"application/json"},
			expectedStatusCode:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := ContentType(tt.allowedContentTypes...)
			handler := middleware(noopHandler)

			req := httptest.NewRequest("POST", "/", strings.NewReader(strings.Repeat("a", int(tt.contentLength))))
			req.Header.Set("Content-Type", tt.contentType)
			req.ContentLength = tt.contentLength

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rec.Code)
			}
		})
	}
}