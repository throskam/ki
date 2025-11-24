package ki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux_Methods(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mux := NewMux()
	mux.Get("/", handler)
	mux.Post("/", handler)
	mux.Put("/", handler)
	mux.Patch("/", handler)
	mux.Delete("/", handler)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected status: got=%v", rec.Code)
			}
			if rec.Body.String() != "ok" {
				t.Fatalf("Unexpected body: got=%q", rec.Body.String())
			}
		})
	}
}

func TestMux_Mount(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	child := NewMux()
	child.Get("/{$}", handler)
	child.Get("/bar", handler)

	parent := NewMux()
	parent.Mount("/foo", child)
	parent.Get("/foo/baz", handler)
	parent.Get("/foobaz", handler)

	tests := []string{
		"/foo/",
		"/foo/bar",
		"/foo/baz",
		"/foobaz",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			parent.ServeHTTP(rec, req)

			t.Logf("Location: %s", rec.Header().Get("Location"))
			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected status for %s: got=%d", path, rec.Code)
			}
			if rec.Body.String() != "ok" {
				t.Fatalf("Unexpected body for %s: got=%q", path, rec.Body.String())
			}
		})
	}
}

func TestMux_Route(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mux := NewMux()
	mux.Route("/foo", func(r Router) {
		r.Get("/{$}", handler)
		r.Get("/bar", handler)
		r.Route("/baz", func(r Router) {
			r.Get("/{$}", handler)
			r.Get("/foobaz", handler)
		})
	})

	tests := []string{
		"/foo/",
		"/foo/bar",
		"/foo/baz/",
		"/foo/baz/foobaz",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected status for %s: got=%d", path, rec.Code)
			}
			if rec.Body.String() != "ok" {
				t.Fatalf("Unexpected body for %s: got=%q", path, rec.Body.String())
			}
		})
	}
}

func TestMux_Group(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mux := NewMux()
	mux.Group(func(r Router) {
		r.Get("/{$}", handler)
		r.Get("/foo", handler)
	})
	mux.Group(func(r Router) {
		r.Get("/bar", handler)
	})
	mux.Route("/baz", func(r Router) {
		r.Get("/{$}", handler)
		r.Group(func(r Router) {
			r.Get("/foobaz", handler)
		})
	})

	paths := []string{
		"/",
		"/foo",
		"/bar",
		"/baz/",
		"/baz/foobaz",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected status for %s: got=%d", path, rec.Code)
			}
			if rec.Body.String() != "ok" {
				t.Fatalf("Unexpected body for %s: got=%q", path, rec.Body.String())
			}
		})
	}
}

func TestMux_Use(t *testing.T) {
	var count int

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mw := func(n int) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count += n
				next.ServeHTTP(w, r)
			})
		}
	}

	mux := NewMux()
	mux.Use(mw(1))
	mux.Get("/foo", handler)

	mux.Route("/bar", func(r Router) {
		r.Use(mw(2))
		r.Get("/baz", handler)

		r.Group(func(r Router) {
			r.Use(mw(3), mw(4))
			r.Get("/foobaz", handler)
		})
	})

	tests := map[string]int{
		"/foo":        1,
		"/bar/baz":    3,
		"/bar/foobaz": 10, // 1 + 2 + 3 + 4
	}

	for path, expected := range tests {
		t.Run(path, func(t *testing.T) {
			count = 0
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected status for %s: got=%d", path, rec.Code)
			}
			if rec.Body.String() != "ok" {
				t.Fatalf("Unexpected body for %s: got=%q", path, rec.Body.String())
			}
			if count != expected {
				t.Fatalf("Wrong middleware count for %s: got=%d, want=%d", path, count, expected)
			}
		})
	}
}

func TestMux_MiddlewareOrder(t *testing.T) {
	var order []int

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mw := func(n int) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, n)
				next.ServeHTTP(w, r)
			})
		}
	}

	mux := NewMux()
	mux.Use(mw(1), mw(2))
	mux.Use(mw(3))
	mux.Get("/", handler, WithMiddleware(mw(4), mw(5)), WithMiddleware(mw(6)))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	order = []int{}
	mux.ServeHTTP(rec, req)

	expected := []int{1, 2, 3, 4, 5, 6}
	for i, got := range order {
		if got != expected[i] {
			t.Fatalf("Wrong middleware order at %d: got=%d, want=%d", i, got, expected[i])
		}
	}
}

func TestMux_Registry(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mux := NewMux()
	mux.Route("/level0", func(r Router) {
		r.Get("/foo", handler, WithName("level0-foo"))
		r.Route("/level1", func(r Router) {
			r.Get("/foo", handler, WithName("level1-foo"))
			r.Route("/level2", func(r Router) {
				r.Get("/foo", handler, WithName("level2-foo"))
			})
		})
	})

	expected := map[string]string{
		"level0-foo": "/level0/foo",
		"level1-foo": "/level0/level1/foo",
		"level2-foo": "/level0/level1/level2/foo",
	}

	for name, path := range expected {
		t.Run(name, func(t *testing.T) {
			loc := mux.Registry().Get(name)
			if loc.Pattern() == "" {
				t.Fatalf("Missing registry entry for %s", name)
			}
			if got := loc.URL().String(); got != path {
				t.Fatalf("Incorrect path: got=%s, want=%s", got, path)
			}
		})
	}
}

func TestMux_404(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d", rec.Code)
	}
}
