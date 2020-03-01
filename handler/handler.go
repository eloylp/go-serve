// Package handler covers all necessary stuff for
// running HTTP server logic.
package handler

import (
	"encoding/base64"
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

// AuthChecker takes as parameter a token for check against  the
// Authorization request header, that needs to be base64 encoded.
// Will let pass the request through the chain if validation
// succeeds. If not, it will stop the chain with an unauthorized
// status code (401) and a "Bad auth" message.
func AuthChecker(token string) www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerAuthB64 := r.Header.Get("Authorization")
			headerAuth, err := base64.StdEncoding.DecodeString(headerAuthB64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Auth header must be encoded in base64"))
				return
			}
			if string(headerAuth) != token {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Bad auth"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
