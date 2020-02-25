package middleware

import (
	"net/http"
)

// Middleware accepts the next http.Handler as parameter
// and returns current one that may modify request/writer
// and finally calls the handler passed as parameter.
type Middleware func(h http.Handler) http.Handler

// Apply will take the handler as first parameter.
// The variadic part of function accepts any number of middlewares
// that will be called in the passed order.
// Beware that the handler will always be called as the
// last element of the chain.
func Apply(h http.Handler, m ...Middleware) http.Handler {
	for j := len(m) - 1; j >= 0; j-- {
		h = m[j](h)
	}
	return h
}
