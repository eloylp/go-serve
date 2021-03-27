package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string) http.Handler {
	r := mux.NewRouter()
	middlewares := []mux.MiddlewareFunc{
		handler.ServerHeader(Version),
		handler.RequestLogger(logger),
	}
	r.Use(middlewares...)
	fileHandler := http.FileServer(http.Dir(docRoot))
	r.Methods(http.MethodGet).PathPrefix(cfg.Prefix).Handler(http.StripPrefix(cfg.Prefix, fileHandler))
	if cfg.UploadEndpoint != "" {
		r.Methods(http.MethodPost).
			Path(cfg.UploadEndpoint).
			Handler(handler.UploadTARGZHandler(logger, cfg.DocRoot))
	}
	return r
}
