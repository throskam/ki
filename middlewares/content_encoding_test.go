package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentEncodingMiddleware(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name               string
		contentEncoding    []string
		contentLength      int64
		expectedStatusCode int
	}{
		{
			name:               "No Content - Should Pass",
			contentEncoding:    nil,
			contentLength:      0,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Allowed Encoding - gzip",
			contentEncoding:    []string{"gzip"},
			contentLength:      100,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Disallowed Encoding - br",
			contentEncoding:    []string{"br"},
			contentLength:      100,
			expectedStatusCode: http.StatusUnsupportedMediaType,
		},
		{
			name:               "Multiple Encodings - One Disallowed",
			contentEncoding:    []string{"gzip", "br"},
			contentLength:      100,
			expectedStatusCode: http.StatusUnsupportedMediaType,
		},
		{
			name:               "Multiple Allowed Encodings",
			contentEncoding:    []string{"gzip", "deflate"},
			contentLength:      100,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowedEncodings := []string{"gzip", "deflate"}

			middleware := ContentEncoding(allowedEncodings...)
			handler := middleware(noopHandler)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(strings.Repeat("x", int(tt.contentLength))))
			req.ContentLength = tt.contentLength

			if tt.contentEncoding != nil {
				for _, encoding := range tt.contentEncoding {
					req.Header.Add("Content-Encoding", encoding)
				}
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rec.Code)
			}
		})
	}
}