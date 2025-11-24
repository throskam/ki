package ki

import (
	"fmt"
	"net/http"
	"slices"
)

// Route represents a route.
type Route struct {
	method      string
	path        string
	handler     http.Handler
	name        string
	middlewares Stack
}

// RouteOption is a function that configures a Route.
type RouteOption func(*Route)

// NewRoute returns a new Route.
func NewRoute(method, path string, handler http.Handler, options ...RouteOption) Route {
	route := Route{
		method:  method,
		path:    path,
		handler: handler,
	}

	for _, o := range options {
		o(&route)
	}

	return route
}

// Name returns the name of the route.
func (r *Route) Name() string {
	return r.name
}

// Method returns the method of the route.
func (r *Route) Method() string {
	return r.method
}

// Path returns the path of the route.
func (r *Route) Path() string {
	return r.path
}

// Pattern returns the pattern of the route.
func (r *Route) Pattern() string {
	if r.method == "" {
		return r.path
	}

	return fmt.Sprintf("%s %s", r.method, r.path)
}

// Handler builds the handler for the route including the middlewares.
func (r *Route) Handler() http.Handler {
	return r.middlewares.Chain(r.handler)
}

// Location returns a new Location for the route.
func (r *Route) Location() Location {
	return NewLocation(r.Method(), r.Path())
}

// WithName returns a new RouteOption that sets the name of the route.
func WithName(name string) RouteOption {
	return func(rc *Route) {
		rc.name = name
	}
}

// WithMiddleware returns a new RouteOption that sets the middlewares for the route.
func WithMiddleware(middlewares ...func(http.Handler) http.Handler) RouteOption {
	slices.Reverse(middlewares)

	return func(rc *Route) {
		rc.middlewares = slices.Concat(middlewares, rc.middlewares)
	}
}
