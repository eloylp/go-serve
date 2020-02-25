package middleware

import (
	"net/http"
)

// Middleware accepts an http.Handler as parameter
// and returns the next
type Middleware func(h http.Handler) http.Handler

// Apply will take the handler as first parameter.
// The variadic function accepts any number of middlewares
// that will be called in the passed order.
// Beware that the handler will always be called as the
// last element of the chain.
func Apply(h http.Handler, m ...Middleware) http.Handler {
	for j := len(m) - 1; j >= 0; j-- {
		h = m[j](h)
	}
	return h
}
