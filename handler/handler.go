// Package handler covers all necessary stuff for
// running HTTP server logic.
package handler

import (
	"fmt"
	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
	"net/http"
)

// ServerHeader will grab server information in the
// "Server" header. Like version.
func ServerHeader(version string) www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Server", fmt.Sprintf("go-serve %s", version))
			h.ServeHTTP(w, r)
		})
	}
}

// RequestLogger will log the client connection
// information on each request.
func RequestLogger(logger logging.Logger) www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s from client %s", r.Method, r.URL.String(), r.RemoteAddr)
			h.ServeHTTP(w, r)
		})
	}
}

// AuthChecker takes the token to check as param and the
// desired response header. Will let pass the request through
// the chain if validation succeeds. If not, will stop the
// chain call by writing a "Bad auth" message.
func AuthChecker(token string, responseCode int) www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != token {
				w.WriteHeader(responseCode)
				_, _ = w.Write([]byte("Bad auth"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
