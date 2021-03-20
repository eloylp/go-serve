// Package handler covers all necessary stuff for
// running HTTP server logic.
package handler

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	auth "github.com/abbot/go-http-auth"
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
func AuthChecker(realm, authFilePath string) mux.MiddlewareFunc {
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

func UploadHandler(docRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployPath := r.Header.Get("GoServe-Deploy-Path")
		uncompressedStream, err := gzip.NewReader(r.Body)
		if err != nil {
			msg := "failed reading compressed gzip: " + err.Error()
			reply(w, http.StatusBadRequest, msg)
			return
		}
		tarReader := tar.NewReader(uncompressedStream)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				msg := "failed reading next part of tar: " + err.Error()
				reply(w, http.StatusInternalServerError, msg)
				return
			}
			switch header.Typeflag {
			case tar.TypeDir:
				path := filepath.Join(docRoot, deployPath, header.Name)
				if err := os.MkdirAll(path, 0755); err != nil {
					msg := "failed reading next part of tar: " + err.Error()
					reply(w, http.StatusInternalServerError, msg)
					return
				}
			case tar.TypeReg:
				path := filepath.Join(docRoot, deployPath, header.Name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					msg := "failed creating dir for next part of tar: " + err.Error()
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				outFile, err := os.Create(path)
				if err != nil {
					msg := "failed creating file part of tar: " + err.Error()
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					msg := "failed copying data part of tar: " + err.Error()
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				_ = outFile.Close()
			default:
				msg := fmt.Sprintf("unknown part of tar: type: %v in %s", header.Typeflag, header.Name)
				reply(w, http.StatusBadRequest, msg)
				return
			}
		}
		reply(w, http.StatusOK, "upload complete !")
	}
}

func reply(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}
