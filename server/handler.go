package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/eloylp/kit/archive"
	"github.com/eloylp/kit/pathutil"
	"github.com/sirupsen/logrus"
)

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
		writtenBytes, err := archive.ExtractTARGZ(r.Body, path)
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
