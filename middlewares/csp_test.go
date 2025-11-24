package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestContentSecurityPolicy(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name     string
		policy   url.Values
		expected string
	}{
		{
			name:     "Empty policy",
			policy:   url.Values{},
			expected: "",
		},
		{
			name: "Simple policy - default-src only",
			policy: url.Values{
				"default-src": {"'self'"},
			},
			expected: "default-src 'self'",
		},
		{
			name: "Common policy with multiple directives",
			policy: url.Values{
				"default-src": {"'self'"},
				"script-src":  {"'self'", "https://cdn.example.com"},
				"style-src":   {"'self'", "https://fonts.googleapis.com"},
				"object-src":  {"'none'"},
			},
			expected: "default-src 'self'; object-src 'none'; script-src 'self' https://cdn.example.com; style-src 'self' https://fonts.googleapis.com",
		},
		{
			name: "Policy with data: and unsafe-inline",
			policy: url.Values{
				"img-src":   {"'self'", "data:"},
				"style-src": {"'self'", "'unsafe-inline'"},
			},
			expected: "img-src 'self' data:; style-src 'self' 'unsafe-inline'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := ContentSecurityPolicy(tt.policy)
			handler := middleware(noopHandler)

			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			got := rec.Header().Get("Content-Security-Policy")

			if got != tt.expected {
				t.Errorf("expected header %q, got %q", tt.expected, got)
			}
		})
	}
}
