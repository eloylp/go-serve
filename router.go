package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
)

func router(cfg *config.Settings, logger *logrus.Logger, serverIdentity string) http.Handler {
	m := mux.NewRouter()
	middlewares := []mux.MiddlewareFunc{
		handler.ServerHeader(Version),
		handler.RequestLogger(logger),
	}
	m.Use(middlewares...)
	if cfg.AuthFile != "" {
		m.Use(handler.AuthChecker(serverIdentity, cfg.AuthFile))
	}
	fileHandler := http.FileServer(http.Dir(cfg.DocRoot))
	m.PathPrefix(cfg.Prefix).Handler(http.StripPrefix(cfg.Prefix, fileHandler))
	return m
}
