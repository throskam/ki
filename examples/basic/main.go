package main

import (
	"net/http"

	"github.com/throskam/ki"
	"github.com/throskam/ki/middlewares"
)

func main() {
	router := ki.NewRouter()

	// Basic middlewares stack.
	router.Use(middlewares.RequestID())
	router.Use(middlewares.RealIP())
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.Recoverer())

	// Add other middlewares here.

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})

	_ = http.ListenAndServe(":8080", router)
}
