package ki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeMiddleware(id string, calls *[]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*calls = append(*calls, "before:"+id)
			next.ServeHTTP(w, r)
			*calls = append(*calls, "after:"+id)
		})
	}
}

func TestStack_Chain(t *testing.T) {
	var calls []string

	m1 := makeMiddleware("m1", &calls)
	m2 := makeMiddleware("m2", &calls)

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	})

	stack := Stack{m1, m2}
	handler := stack.Chain(noopHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := []string{
		"before:m2",
		"before:m1",
		"handler",
		"after:m1",
		"after:m2",
	}

	if len(calls) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(calls))
	}

	for i, call := range expected {
		if calls[i] != call {
			t.Errorf("expected call %d to be %q, got %q", i, call, calls[i])
		}
	}
}
