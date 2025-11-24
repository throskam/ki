package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealIPMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		expectedRemote string
	}{
		{
			name: "True-Client-IP set",
			headers: map[string]string{
				"True-Client-IP": "1.2.3.4",
			},
			expectedRemote: "1.2.3.4",
		},
		{
			name: "X-Real-IP set",
			headers: map[string]string{
				"X-Real-IP": "5.6.7.8",
			},
			expectedRemote: "5.6.7.8",
		},
		{
			name: "X-Forwarded-For set with multiple IPs",
			headers: map[string]string{
				"X-Forwarded-For": "9.10.11.12, 13.14.15.16",
			},
			expectedRemote: "9.10.11.12",
		},
		{
			name:           "No headers, fallback to original RemoteAddr",
			headers:        map[string]string{},
			expectedRemote: "192.168.0.1:1234",
		},
		{
			name: "Invalid IP in headers, fallback to original RemoteAddr",
			headers: map[string]string{
				"X-Real-IP": "invalid-ip",
			},
			expectedRemote: "192.168.0.1:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RemoteAddr != tt.expectedRemote {
					t.Errorf("expected RemoteAddr to be %q, got %q", tt.expectedRemote, r.RemoteAddr)
				}
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			req.RemoteAddr = "192.168.0.1:1234"

			rec := httptest.NewRecorder()

			middleware := RealIP()
			handler := middleware(noopHandler)

			handler.ServeHTTP(rec, req)
		})
	}
}