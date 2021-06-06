package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"go.eloylp.dev/go-serve/config"
)

var UploadSize *prometheus.HistogramVec

func uploadSize(buckets []float64) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "goserve",
		Subsystem: "upload",
		Name:      "size",
		Help:      "Histogram to represent the successful uploads to the server",
		Buckets:   buckets,
	}, []string{})
}

func Initialize(cfg *config.Settings) {
	UploadSize = uploadSize(cfg.MetricsSizeBuckets)
	prometheus.MustRegister(UploadSize)
}
