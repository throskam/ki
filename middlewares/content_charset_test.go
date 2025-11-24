package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentCharset(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name            string
		contentType     string
		contentLength   int64
		allowedCharsets []string
		expectedStatus  int
	}{
		{
			name:            "Allowed charset utf-8",
			contentType:     "application/json; charset=utf-8",
			contentLength:   10,
			allowedCharsets: []string{"utf-8"},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Allowed charset iso-8859-1",
			contentType:     "text/plain; charset=iso-8859-1",
			contentLength:   5,
			allowedCharsets: []string{"utf-8", "iso-8859-1"},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Unsupported charset",
			contentType:     "application/json; charset=windows-1252",
			contentLength:   20,
			allowedCharsets: []string{"utf-8"},
			expectedStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:            "No charset provided, default utf-8 allowed",
			contentType:     "application/json",
			contentLength:   15,
			allowedCharsets: []string{"utf-8"},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "No charset provided, default utf-8 not allowed",
			contentType:     "application/json",
			contentLength:   15,
			allowedCharsets: []string{"iso-8859-1"},
			expectedStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:            "Zero Content-Length should pass",
			contentType:     "application/json; charset=windows-1252",
			contentLength:   0,
			allowedCharsets: []string{},
			expectedStatus:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := ContentCharset(tt.allowedCharsets...)
			handler := middleware(noopHandler)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(strings.Repeat("a", int(tt.contentLength))))
			req.Header.Set("Content-Type", tt.contentType)
			req.ContentLength = tt.contentLength

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
