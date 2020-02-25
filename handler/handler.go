package handler

import (
	"fmt"
	"github.com/eloylp/go-serve/www"
	"log"
	"net/http"
)

func VersionHeader(version string) www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Server", fmt.Sprintf("go-serve %s", version))
			h.ServeHTTP(w, r)
		})
	}
}

func RequestLogger() www.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			go log.Printf("%s %s from client %s", r.Method, r.RequestURI, r.RemoteAddr)
			h.ServeHTTP(w, r)
		})
	}
}
