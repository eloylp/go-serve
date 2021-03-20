// Package config is a local use package for grouping
// app config related functions
package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ListenAddr      string `default:"0.0.0.0:8080"`
	DocRoot         string `required:"."`
	Prefix          string `default:"/"`
	UploadEndpoint  string
	ShutdownTimeout time.Duration `default:"1s"`
	Logger          *LoggerSettings
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"25s"`
}

type LoggerSettings struct {
	Format string `default:"json"`
	Output io.Writer
	Level  string `default:"info"`
}

func FromEnv() (*Settings, error) {
	s := defaultSettings()
	if err := envconfig.Process("GOSERVE", s); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return s, nil
}

func defaultSettings() *Settings {
	s := &Settings{
		ListenAddr:      "0.0.0.0:8080",
		DocRoot:         ".",
		Prefix:          "/",
		ShutdownTimeout: time.Second,
		Logger: &LoggerSettings{
			Level:  logrus.InfoLevel.String(),
			Format: "text",
			Output: os.Stderr,
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 25 * time.Second,
	}
	return s
}
