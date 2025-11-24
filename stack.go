package ki

import "net/http"

// Stack is a stack of middlewares.
type Stack []func(http.Handler) http.Handler

// Chain returns the chained handler.
func (s *Stack) Chain(handler http.Handler) http.Handler {
	chain := handler

	for _, m := range *s {
		chain = m(chain)
	}

	return chain
}