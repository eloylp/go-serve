// Package config is a local use package for grouping
// app config related functions
package config

import (
	"fmt"
	"io"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	ListenAddr      string `default:"0.0.0.0:8080"`
	DocRoot         string `required:"."`
	Prefix          string `default:"/"`
	AuthFile        string
	ShutdownTimeout time.Duration `default:"1s"`
	Logger          *LoggerSettings
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"25s"`
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
