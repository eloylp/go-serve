// Package handler covers all necessary stuff for
// running HTTP server logic.
package handler

import (
	"fmt"
	"github.com/eloylp/go-serve/packer"
	"net/http"
	"path/filepath"

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

func UploadTARGZHandler(logger *logrus.Logger, docRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployPath := r.Header.Get("GoServe-Deploy-Path")
		writtenBytes, err := ProcessTARGZStream(r.Body, docRoot, deployPath)
		if err != nil {
			logger.Debugf("%v", err)
			reply(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := fmt.Sprintf("upload of tar.gz complete ! Bytes written: %d", writtenBytes)
		logger.Debug(msg)
		reply(w, http.StatusOK, msg)
	}
}

func DownloadTARGZHandler(logger *logrus.Logger, root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		downloadRelativePath := r.Header.Get("GoServe-Download-Path")
		downloadAbsolutePath := filepath.Join(root, downloadRelativePath)
		if err := checkPath(root, downloadAbsolutePath); err != nil {
			logger.WithError(err).Error("download path violation try")
			reply(w, http.StatusBadRequest, err.Error())
			return
		}
		writtenBytes, err := packer.WriteTARGZ(w, downloadAbsolutePath)
		if err != nil {
			logger.WithError(err).Error("fail writing tar.gz to wire")
			return
		}
		logger.Debugf("sent of tar.gz to %s complete ! Bytes written: %d", r.RemoteAddr, writtenBytes)
	}
}

func reply(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}
