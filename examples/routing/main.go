package main

import (
	"fmt"
	"net/http"

	"github.com/throskam/ki"
)

func main() {
	router := ki.NewRouter()

	// Add a middleware to the router.
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("middleware\n"))
			next.ServeHTTP(w, r)
		})
	})

	// GET /
	router.Get("/{$}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("home\n"))
	})

	router.Route("/posts", func(r ki.Router) {
		// GET /posts
		r.Get("/{$}", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("list posts\n"))
		})

		// POST /posts
		r.Post("/{$}", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("create a new post\n"))
		})

		r.Group(func(g ki.Router) {
			// Add a middleware to the group.
			g.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintf(w, "post ID: %s\n", r.PathValue("id"))
					next.ServeHTTP(w, r)
				})
			})

			// GET /posts/{id}
			g.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("get a post\n"))
			})

			// DELETE /posts/{id}
			g.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("delete a post\n"))
			})
		})
	})

	// Create an admin router.
	admin := ki.NewRouter()

	// Add a middleware to the admin router.
	admin.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("admin\n"))
			next.ServeHTTP(w, r)
		})
	})

	// GET /admin/dashboard
	admin.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("dashboard\n"))
	})

	// Mount the admin router at the given prefix.
	router.Mount("/admin", admin)

	_ = http.ListenAndServe(":8080", router)
}
