// Package config is a local use package for grouping
// app config related functions
package config

import (
	"flag"
	"os"
	"time"
)

type Option func(cfg *Settings)

type Settings struct {
	DocRoot, Prefix, ListenAddr, AuthFile string
	ShutdownTimeout                       time.Duration
}

// FromArgs will receive and argument list as parameter, including
// the program name and returning the proper variables with all the values.
func FromArgs(args []string) (s Settings, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	flag.StringVar(&s.DocRoot, "d", currentDir, "Defines the document root")
	flag.StringVar(&s.Prefix, "p", "/", "Defines prefix to use for serve files")
	flag.StringVar(&s.ListenAddr, "l", "0.0.0.0:8080", "Defines the listen address")
	flag.StringVar(&s.AuthFile, "a", "", "Defines the .htpasswd file path for auth")
	argsFiltered := args[1:] // exclude program name
	err = flag.CommandLine.Parse(argsFiltered)
	return
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
	cfg := &Settings{}
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
