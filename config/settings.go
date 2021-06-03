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
	ListenAddr                    string          `default:"0.0.0.0:8080" split_words:"true"`
	DocRoot                       string          `required:"." split_words:"true"`
	Prefix                        string          `default:"/static" split_words:"true"`
	UploadEndpoint                string          `split_words:"true"`
	DownloadEndpoint              string          `split_words:"true"`
	ShutdownTimeout               time.Duration   `default:"5s" split_words:"true"`
	Logger                        *LoggerSettings `split_words:"true"`
	ReadTimeout                   time.Duration   `default:"0s" split_words:"true"`
	WriteTimeout                  time.Duration   `default:"0s" split_words:"true"`
	ReadAuthorizations            Authorization   `split_words:"true"`
	WriteAuthorizations           Authorization   `split_words:"true"`
	MetricsEnabled                bool            `default:"true" split_words:"true"`
	MetricsPath                   string          `default:"/metrics" split_words:"true"`
	MetricsListenAddr             string          `split_words:"true"`
	MetricsRequestDurationBuckets []float64       `split_words:"true"`
	MetricsSizeBuckets            []float64       `split_words:"true"`
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
		Prefix:          "/static",
		ShutdownTimeout: time.Second,
		Logger: &LoggerSettings{
			Level:  logrus.InfoLevel.String(),
			Format: "text",
			Output: os.Stderr,
		},
		ReadTimeout:                   0,
		WriteTimeout:                  0,
		WriteAuthorizations:           Authorization{},
		ReadAuthorizations:            Authorization{},
		MetricsEnabled:                true,
		MetricsPath:                   "/metrics",
		MetricsRequestDurationBuckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		MetricsSizeBuckets:            []float64{1_000, 10_000, 50_000, 100_000, 500_000, 1_000_000, 2_000_000, 4_000_000, 8_000_000},
	}
	return s
}
