// Package handler covers all necessary stuff for
// running HTTP server logic.
package handler

import (
	"fmt"
	"net/http"

	auth "github.com/abbot/go-http-auth"

	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
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

// AuthChecker represents the basic auth implementation
// https://tools.ietf.org/html/rfc7617
// Will let pass the request through the chain if validation
// succeeds. If not, it will stop the chain with an unauthorized
// status code (401).
// A status code (500) with  "Bad auth file" message as body
// will be returned if the basic auth file is not correct.
// The underlying library will watch the file for changes
// and will update the server automatically.
func AuthChecker(realm, authFilePath string) www.Middleware {
	ap := auth.HtpasswdFileProvider(authFilePath)
	authenticator := auth.NewBasicAuthenticator(realm, ap)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("Bad auth file"))
					return
				}
			}()
			if authenticator.CheckAuth(r) == "" {
				authenticator.RequireAuth(w, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}