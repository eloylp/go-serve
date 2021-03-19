package server

import (
	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
)

func loggerFrom(cfg *config.LoggerSettings) *logrus.Logger {
	l := logrus.New()
	if cfg.Output != nil {
		l.SetOutput(cfg.Output)
	}
	if cfg.Format == "json" {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetFormatter(&logrus.TextFormatter{})
	}
	return l
}
