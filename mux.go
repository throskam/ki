package ki

import (
	"fmt"
	"net/http"
	"slices"
)

// Mux is a router that uses a ServeMux.
type Mux struct {
	mux *http.ServeMux

	registry *Registry

	routeOptions []RouteOption
}

// NewMux returns a new Mux.
func NewMux() *Mux {
	return &Mux{
		mux:          http.NewServeMux(),
		registry:     NewRegistry(),
		routeOptions: []RouteOption{},
	}
}

// ServeHTTP implements the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

// Mount mounts the given handler at the given prefix.
func (m *Mux) Mount(prefix string, handler http.Handler) {
	m.mount(prefix, handler)
}

// Route creates a new router with the given prefix.
func (m *Mux) Route(prefix string, fn func(Router)) Router {
	mux := &Mux{
		mux:          http.NewServeMux(),
		registry:     m.registry.Child(prefix),
		routeOptions: []RouteOption{},
	}

	m.Mount(prefix, mux)

	if fn != nil {
		fn(mux)
	}

	return mux
}

// Group creates a new router without any prefix.
func (m *Mux) Group(fn func(Router)) Router {
	mux := &Mux{
		mux:          m.mux,
		registry:     m.registry,
		routeOptions: slices.Clone(m.routeOptions),
	}

	if fn != nil {
		fn(mux)
	}

	return mux
}

// Use adds the given middlewares to the router.
func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.routeOptions = append(m.routeOptions, WithMiddleware(middlewares...))
}

// Method adds a route for the given verb.
func (m *Mux) Method(method, pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.method(method, pattern, handler, options...)
}

// Get adds a route for the GET verb.
func (m *Mux) Get(pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.Method(http.MethodGet, pattern, handler, options...)
}

// Post adds a route for the POST verb.
func (m *Mux) Post(pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.Method(http.MethodPost, pattern, handler, options...)
}

// Put adds a route for the PUT verb.
func (m *Mux) Put(pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.Method(http.MethodPut, pattern, handler, options...)
}

// Patch adds a route for the PATCH verb.
func (m *Mux) Patch(pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.Method(http.MethodPatch, pattern, handler, options...)
}

// Delete adds a route for the DELETE verb.
func (m *Mux) Delete(pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	return m.Method(http.MethodDelete, pattern, handler, options...)
}

// Registry returns the registry of the router.
func (m *Mux) Registry() *Registry {
	return m.registry
}

// handle adds the route to the mux.
func (m *Mux) handle(route Route) {
	m.mux.Handle(route.Pattern(), route.Handler())
}

// mount mounts the given handler at the given prefix.
func (m *Mux) mount(prefix string, handler http.Handler) {
	pattern := fmt.Sprintf("%s/", prefix)

	route := NewRoute("", pattern, handler, slices.Concat(m.routeOptions, []RouteOption{WithMiddleware(
		func(next http.Handler) http.Handler { return http.StripPrefix(prefix, next) },
	)})...)

	m.handle(route)
}

// method adds a route for the given verb.
func (m *Mux) method(method, pattern string, handler http.HandlerFunc, options ...RouteOption) Location {
	route := NewRoute(method, pattern, handler, slices.Concat(m.routeOptions, options)...)

	if route.Name() != "" {
		m.registry.Add(route.Name(), route.Method(), route.Path())
	}

	m.handle(route)

	return route.Location()
}