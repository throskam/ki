package ki

import (
	"fmt"
)

// Registry is a registry of routes.
type Registry struct {
	routeMap   map[string]Location
	registries map[string]*Registry
}

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		routeMap:   map[string]Location{},
		registries: map[string]*Registry{},
	}
}

// Add adds a route to the registry.
// It panics if the route already exists.
func (r *Registry) Add(key, method, pattern string) {
	_, ok := r.routeMap[key]

	if ok {
		panic(fmt.Sprintf("Location %s already exists", key))
	}

	r.routeMap[key] = NewLocation(method, pattern)
}

// Remove removes a route from the registry.
func (r *Registry) Remove(key string) {
	delete(r.routeMap, key)
}

// Has returns true if the registry has a route with the given key or any of its child registries.
func (r *Registry) Has(key string) bool {
	_, ok := r.routeMap[key]

	if !ok {
		for _, registry := range r.registries {
			if registry.Has(key) {
				return true
			}
		}

		return false
	}

	return true
}

// Get returns the location for the given key.
// It panics if the location does not exist.
func (r *Registry) Get(key string) Location {
	location, ok := r.routeMap[key]

	if !ok {
		for prefix, registry := range r.registries {
			if registry.Has(key) {
				return registry.Get(key).WithPrefix(prefix)
			}
		}

		panic(fmt.Sprintf("Location %s does not exist", key))
	}

	return location
}

// Child returns a new child registry with the given prefix.
func (r *Registry) Child(prefix string) *Registry {
	registry := NewRegistry()

	r.registries[prefix] = registry

	return registry
}
