package config

import (
	"io"
	"time"
)

type Option func(cfg *Settings)

func ForOptions(opts ...Option) *Settings {
	cfg := defaultSettings()
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

func WithReadTimeout(duration time.Duration) Option {
	return func(cfg *Settings) {
		cfg.ReadTimeout = duration
	}
}

func WithWriteTimeout(duration time.Duration) Option {
	return func(cfg *Settings) {
		cfg.WriteTimeout = duration
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

func WithUploadEndpoint(path string) Option {
	return func(cfg *Settings) {
		cfg.UploadEndpoint = path
	}
}

func WithDownLoadEndpoint(path string) Option {
	return func(cfg *Settings) {
		cfg.DownloadEndpoint = path
	}
}

func WithLoggerLevel(level string) Option {
	return func(cfg *Settings) {
		cfg.Logger.Level = level
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
