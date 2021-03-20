package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string, lock *sync.RWMutex) http.Handler {
	r := mux.NewRouter()
	middlewares := []mux.MiddlewareFunc{
		handler.ServerHeader(Version),
		handler.RequestLogger(logger),
	}
	r.Use(middlewares...)
	fileHandler := http.FileServer(http.Dir(docRoot))
	syncFileHandler := handler.SyncRead(http.StripPrefix(cfg.Prefix, fileHandler), lock)
	r.Methods(http.MethodGet).PathPrefix(cfg.Prefix).Handler(syncFileHandler)
	if cfg.UploadEndpoint != "" {
		syncUploadTARGZHandler := handler.SyncWrite(handler.UploadTARGZHandler(logger, cfg.DocRoot), lock)
		r.Methods(http.MethodPost).
			Path(cfg.UploadEndpoint).
			Handler(syncUploadTARGZHandler)
	}
	return r
}
