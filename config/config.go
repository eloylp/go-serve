// Package config is a local use package for grouping
// app config related functions
package config

import (
	"fmt"
	"io"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Option func(cfg *Settings)

type Settings struct {
	ListenAddr      string `default:"0.0.0.0:8080"`
	DocRoot         string `default:"."`
	Prefix          string `default:"/"`
	AuthFile        string
	ShutdownTimeout time.Duration `default:"1s"`
	Logger          *LoggerSettings
}

type LoggerSettings struct {
	Format string `default:"json"`
	Output io.Writer
}

func FromEnv() (*Settings, error) {
	s := emptySettings()
	if err := envconfig.Process("GOSERVE", s); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return s, nil
}

func emptySettings() *Settings {
	s := &Settings{
		Logger: &LoggerSettings{},
	}
	return s
}

func WithListenAddr(addr string) Option {
	return func(cfg *Settings) {
		cfg.ListenAddr = addr
	}
}

func WithDocRoot(docRoot string) Option {
	return func(cfg *Settings) {
		cfg.DocRoot = docRoot
	}
}

func ForOptions(opts ...Option) *Settings {
	cfg := emptySettings()
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func WithDocRootPrefix(prefix string) Option {
	return func(cfg *Settings) {
		cfg.Prefix = prefix
	}
}

func WithLoggerFormat(format string) Option {
	return func(cfg *Settings) {
		cfg.Logger.Format = format
	}
}

func WithLoggerOutput(o io.Writer) Option {
	return func(cfg *Settings) {
		cfg.Logger.Output = o
	}
}
