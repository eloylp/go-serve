package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.eloylp.dev/kit/archive"
	"go.eloylp.dev/kit/pathutil"

	"github.com/eloylp/go-serve/metrics"
)

const ContentTypeTarGzip = "application/tar+gzip"
const ContentTypeFile = "application/octet-stream"

func StatusHandler(info Info) http.HandlerFunc {
	type Status struct {
		Status string `json:"status"`
		Info   Info   `json:"info"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(&Status{
			Status: "ok",
			Info:   info,
		})
	}
}

func UploadTARGZHandler(logger *logrus.Logger, docRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deployPath := r.Header.Get("GoServe-Deploy-Path")
		path := filepath.Join(docRoot, deployPath) // nolinter: gosec
		if err := pathutil.PathInRoot(docRoot, path); err != nil {
			logger.WithError(err).Error("upload path violation try")
			reply(w, http.StatusBadRequest, err.Error())
			return
		}
		var writtenBytes int64
		var err error
		switch r.Header.Get("Content-Type") {
		case ContentTypeTarGzip:
			writtenBytes, err = archive.ExtractTARGZ(r.Body, path)
			if err != nil {
				logger.Debugf("%v", err)
				reply(w, http.StatusBadRequest, err.Error())
				return
			}
		case ContentTypeFile:
			writtenBytes, err = saveFile(r.Body, path)
			if err != nil {
				logger.Debugf("%v", err)
				reply(w, http.StatusBadRequest, err.Error())
				return
			}
		default:
			http.NotFound(w, r)
			return
		}
		msg := fmt.Sprintf("upload complete ! Bytes written: %d", writtenBytes)
		logger.Debug(msg)
		if metrics.UploadSize != nil {
			metrics.UploadSize.WithLabelValues().Observe(float64(writtenBytes))
		}
		reply(w, http.StatusOK, msg)
	}
}

func saveFile(reader io.Reader, path string) (int64, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, err
	}
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(file, reader)
	if err != nil {
		return 0, err
	}
	return written, nil
}

func DownloadTARGZHandler(logger *logrus.Logger, root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != ContentTypeTarGzip {
			http.NotFound(w, r)
			return
		}
		downloadRelativePath := r.Header.Get("GoServe-Download-Path")
		downloadAbsolutePath := filepath.Join(root, downloadRelativePath)
		if err := pathutil.PathInRoot(root, downloadAbsolutePath); err != nil {
			logger.WithError(err).Error("download path violation try")
			reply(w, http.StatusBadRequest, err.Error())
			return
		}
		writtenBytes, err := archive.CreateTARGZ(w, downloadAbsolutePath)
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
