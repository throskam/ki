package middlewares

import (
	"net/http"

	"github.com/throskam/ki"
	"golang.org/x/text/language"
)

// Language returns a middleware that sets the preferred language for the request.
func Language(supportedLanguages ...language.Tag) func(http.Handler) http.Handler {
	matcher := language.NewMatcher(supportedLanguages)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := getPreferredLanguage(r, matcher)

			ctx := ki.SetLanguage(r.Context(), lang)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OverrideLanguage returns a middleware that overrides the preferred language for the request based on the value of the cookie.
func OverrideLanguage(cookie string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang, err := r.Cookie(cookie)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if lang == nil || lang.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			tag := language.MustParse(lang.Value)

			ctx := ki.SetLanguage(r.Context(), tag)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getPreferredLanguage returns the preferred language for the request based on the Accept-Language header.
func getPreferredLanguage(r *http.Request, matcher language.Matcher) language.Tag {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return language.English
	}

	tag, _, _ := matcher.Match(tags...)

	return tag
}

