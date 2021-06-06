package server

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"go.eloylp.dev/go-serve/config"
)

func logger(cfg *config.LoggerSettings) (*logrus.Logger, error) {
	l := logrus.New()
	l.SetOutput(cfg.Output)
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	l.SetLevel(level)
	if cfg.Format == "json" {
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		l.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}
	return l, nil
}
