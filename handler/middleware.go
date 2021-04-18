package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// ServerHeader will grab server information in the
// "Server" header. Like version.
func ServerHeader(version string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Server", fmt.Sprintf("go-serve %s", version))
			h.ServeHTTP(w, r)
		})
	}
}

// RequestLogger will log the client connection
// information on each request.
func RequestLogger(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(logrus.Fields{
				"path":    r.URL.String(),
				"method":  r.Method,
				"ip":      r.RemoteAddr,
				"headers": r.Header,
			}).Info("request from client")
			h.ServeHTTP(w, r)
		})
	}
}
