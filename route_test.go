package ki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoute_Basic(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	route := NewRoute("GET", "/example", noopHandler)

	if route.Method() != "GET" {
		t.Errorf("expected method 'GET', got '%s'", route.Method())
	}

	if route.Path() != "/example" {
		t.Errorf("expected path '/example', got '%s'", route.Path())
	}

	if route.Pattern() != "GET /example" {
		t.Errorf("expected pattern 'GET /example', got '%s'", route.Pattern())
	}

	if route.Name() != "" {
		t.Errorf("expected empty name, got '%s'", route.Name())
	}
}

func TestRoute_WithNameOption(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	route := NewRoute("POST", "/named", noopHandler, WithName("MyRoute"))

	if route.Name() != "MyRoute" {
		t.Errorf("expected name 'MyRoute', got '%s'", route.Name())
	}
}

func TestRoute_PatternWithNoMethod(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	route := NewRoute("", "/onlypath", noopHandler)

	if route.Pattern() != "/onlypath" {
		t.Errorf("expected pattern '/onlypath', got '%s'", route.Pattern())
	}
}

func TestRoute_WithMiddlewareOrder(t *testing.T) {
	callOrder := []string{}

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw1")
			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw2")
			next.ServeHTTP(w, r)
		})
	}

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	route := NewRoute("GET", "/middleware", noopHandler, WithMiddleware(mw1, mw2))

	req := httptest.NewRequest("GET", "/middleware", nil)
	rec := httptest.NewRecorder()

	route.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rec.Code)
	}

	expectedOrder := []string{"mw1", "mw2", "handler"}

	if len(callOrder) != len(expectedOrder) {
		t.Errorf("expected %d calls, got %d", len(expectedOrder), len(callOrder))
	}

	for i := range expectedOrder {
		if callOrder[i] != expectedOrder[i] {
			t.Errorf("at index %d: expected '%s', got '%s'", i, expectedOrder[i], callOrder[i])
		}
	}
}

func TestRoute_LocationCreation(t *testing.T) {
	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	route := NewRoute("PUT", "/loc", noopHandler)

	loc := route.Location()

	if loc.Method() != "PUT" {
		t.Errorf("expected Location method 'PUT', got '%s'", loc.Method())
	}

	if loc.Pattern() != "/loc" {
		t.Errorf("expected Location pattern '/loc', got '%s'", loc.Pattern())
	}
}
