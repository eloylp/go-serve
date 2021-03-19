package config

import (
	"io"
)

type Option func(cfg *Settings)

func ForOptions(opts ...Option) *Settings {
	cfg := emptySettings()
	for _, o := range opts {
		o(cfg)
	}
	return cfg
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
