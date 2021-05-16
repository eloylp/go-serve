package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ListenAddr                    string `default:"0.0.0.0:8080"`
	DocRoot                       string `required:"."`
	Prefix                        string `default:"/"`
	UploadEndpoint                string
	DownloadEndpoint              string
	ShutdownTimeout               time.Duration `default:"1s"`
	Logger                        *LoggerSettings
	ReadTimeout                   time.Duration `default:"0s"`
	WriteTimeout                  time.Duration `default:"0s"`
	ReadAuthorizations            Authorization
	WriteAuthorizations           Authorization
	PrometheusRegistry            *prometheus.Registry
	MetricsEnabled                bool   `default:"true"`
	MetricsPath                   string `default:"/metrics"`
	MetricsAlternativeListenAddr  string
	MetricsRequestDurationBuckets []float64
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
		ReadTimeout:         0,
		WriteTimeout:        0,
		WriteAuthorizations: Authorization{},
		ReadAuthorizations:  Authorization{},
		MetricsEnabled:      true,
		MetricsPath:         "/metrics",
		PrometheusRegistry:  prometheus.NewRegistry(),
	}
	return s
}
