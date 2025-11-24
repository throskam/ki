package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/throskam/ki"
	"golang.org/x/text/language"
)

func TestLanguageMiddleware(t *testing.T) {
	supported := []language.Tag{language.English, language.Spanish, language.French}
	middleware := Language(supported...)

	tests := []struct {
		name           string
		acceptLanguage string
		expectedLang   language.Tag
	}{
		{
			name:           "Exact match",
			acceptLanguage: "es",
			expectedLang:   language.Spanish,
		},
		{
			name:           "Multiple with priority",
			acceptLanguage: "fr;q=0.9,en;q=0.5",
			expectedLang:   language.French,
		},
		{
			name:           "Unsupported fallback",
			acceptLanguage: "de",
			expectedLang:   language.English,
		},
		{
			name:           "Empty header",
			acceptLanguage: "",
			expectedLang:   language.English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			req.Header.Set("Accept-Language", tt.acceptLanguage)

			noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lang := ki.MustGetLanguage(r.Context())
				if lang != tt.expectedLang {
					t.Errorf("expected language %q, got %q", tt.expectedLang, lang)
				}

				w.WriteHeader(http.StatusOK)
			})

			handler := middleware(noopHandler)

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
		})
	}
}

func TestOverrideLanguageMiddleware(t *testing.T) {
	cookieName := "lang"
	middleware := OverrideLanguage(cookieName)

	tests := []struct {
		name         string
		cookieValue  string
		expectedLang language.Tag
		setCookie    bool
	}{
		{
			name:         "Valid cookie",
			cookieValue:  "fr",
			expectedLang: language.French,
			setCookie:    true,
		},
		{
			name:         "Empty cookie value",
			cookieValue:  "",
			expectedLang: language.English,
			setCookie:    true,
		},
		{
			name:         "No cookie",
			expectedLang: language.English,
			setCookie:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(ki.SetLanguage(req.Context(), language.English))

			if tt.setCookie {
				req.AddCookie(&http.Cookie{
					Name:  cookieName,
					Value: tt.cookieValue,
				})
			}

			rec := httptest.NewRecorder()

			noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lang := ki.MustGetLanguage(r.Context())
				if lang != tt.expectedLang {
					t.Errorf("expected language %q, got %q", tt.expectedLang, lang)
				}

				w.WriteHeader(http.StatusOK)
			})

			handler := middleware(noopHandler)

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
		})
	}
}

