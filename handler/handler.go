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
	"strings"

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

func UploadTARGZHandler(logger *logrus.Logger, docRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployPath := r.Header.Get("GoServe-Deploy-Path")
		uncompressedStream, err := gzip.NewReader(r.Body)
		if err != nil {
			msg := "failed reading compressed gzip: " + err.Error()
			logger.Error(msg)
			reply(w, http.StatusBadRequest, msg)
			return
		}
		var writtenBytes int64
		tarReader := tar.NewReader(uncompressedStream)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				msg := "failed reading next part of tar: " + err.Error()
				logger.Error(msg)
				reply(w, http.StatusInternalServerError, msg)
				return
			}
			// Check that path does not go outside the document root
			path := filepath.Join(docRoot, deployPath, header.Name) // nolinter: gosec
			if err := checkPath(docRoot, path); err != nil {
				msg := "incorrect upload path: " + err.Error()
				logger.Debug(msg)
				reply(w, http.StatusBadRequest, msg)
				return
			}
			// Start processing types
			switch header.Typeflag {
			case tar.TypeDir:
				if err := os.MkdirAll(path, 0755); err != nil {
					msg := fmt.Sprintf("failed creating dir %s part of tar: "+err.Error(), path)
					logger.Error(msg)
					reply(w, http.StatusInternalServerError, msg)
					return
				}
			case tar.TypeReg:
				dir := filepath.Dir(path)
				if err := os.MkdirAll(dir, 0755); err != nil {
					msg := fmt.Sprintf("failed creating dir %s part of tar: "+err.Error(), dir)
					logger.Error(msg)
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				outFile, err := os.Create(path)
				if err != nil {
					msg := fmt.Sprintf("failed creating file part of tar: "+err.Error(), path)
					logger.Error(msg)
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				fileBytes, err := io.Copy(outFile, tarReader) // nolinter: gosec (controlled by read/write timeouts)
				if err != nil {
					msg := fmt.Sprintf("failed copying data of file %s part of tar: %v", path, err)
					logger.Error(msg)
					reply(w, http.StatusInternalServerError, msg)
					return
				}
				writtenBytes += fileBytes
				_ = outFile.Close()
			default:
				msg := fmt.Sprintf("unknown part of tar: type: %v in %s", header.Typeflag, header.Name)
				logger.Error(msg)
				reply(w, http.StatusBadRequest, msg)
				return
			}
		}
		msg := fmt.Sprintf("upload of tar.gz complete ! Bytes written: %d", writtenBytes)
		logger.Debug(msg)
		reply(w, http.StatusOK, msg)
	}
}

func checkPath(docRoot, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(abs, docRoot) {
		return fmt.Errorf("the path you provided %s is not a suitable one", path)
	}
	return nil
}

func reply(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}
