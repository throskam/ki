// Package ki is a thin wrapper around the standard Go mux.
package ki

import (
	"net/http"
)

// Router is a router.
// It provides a way to mount handlers, define routes, and use middlewares.
type Router interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	// Mount mounts the given handler at the given prefix.
	Mount(prefix string, handler http.Handler)

	// Route creates a new router with the given prefix.
	Route(prefix string, fn func(Router)) Router

	// Group creates a new router without any prefix.
	// It is useful for adding middlewares to a group of routes.
	Group(fn func(Router)) Router

	// Use adds the given middlewares to the router.
	Use(middlewares ...func(http.Handler) http.Handler)

	// Method adds a route for the given verb.
	Method(method, pattern string, handler http.HandlerFunc, options ...RouteOption) Location

	// HTTP routing methods.
	Get(pattern string, handler http.HandlerFunc, options ...RouteOption) Location
	Post(pattern string, handler http.HandlerFunc, options ...RouteOption) Location
	Put(pattern string, handler http.HandlerFunc, options ...RouteOption) Location
	Patch(pattern string, handler http.HandlerFunc, options ...RouteOption) Location
	Delete(pattern string, handler http.HandlerFunc, options ...RouteOption) Location

	// Registry returns the registry of the router.
	Registry() *Registry
}

// NewRouter returns a new Router.
func NewRouter() *Mux {
	return NewMux()
}
