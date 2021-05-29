//+build integration

package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestMetricsAreServedByDefault(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	resp, err := http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMetricsAreObserving(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	// Make a request.
	resp, err := http.Get(HTTPAddressStatic)
	assert.NoError(t, err)
	defer resp.Body.Close()

	resp, err = http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	metrics := string(body)
	assert.Contains(t, metrics, "goserve_http_request_duration_seconds_count 1")
}

func TestMetricsCanBeDisabled(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithMetricsEnabled(false),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	AssertDefaultMetricsPathConfigured(t)
}

func AssertDefaultMetricsPathConfigured(t *testing.T) {
	resp, err := http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMetricsCanBeServedOnAlternativePort(t *testing.T) {
	loggerOutput := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(loggerOutput),
			config.WithMetricsAlternativeListenAddr("0.0.0.0:9091"),
		),
	)
	assert.NoError(t, err)
	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)
	test.WaitTCPService(t, "localhost:9091", time.Millisecond, time.Second)

	AssertDefaultMetricsPathConfigured(t)

	resp, err := http.Get("http://localhost:9091/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, loggerOutput.String(), "starting to serve metrics at 0.0.0.0:9091")
}

func TestMetricsRequestDurationBucketsConfig(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithMetricsRequestDurationBuckets([]float64{0.1, 0.5, 1}),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	resp, err := http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	metrics := string(data)

	assert.Contains(t, metrics, "goserve_http_request_duration_seconds_bucket{le=\"0.1\"} 0")
	assert.Contains(t, metrics, "goserve_http_request_duration_seconds_bucket{le=\"0.5\"} 0")
	assert.Contains(t, metrics, "goserve_http_request_duration_seconds_bucket{le=\"1\"} 0")
}

func TestMetricsCanBeServedAlternativePath(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithMetricsPath("/metrics-alternate"),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	resp, err := http.Get(HTTPAddress + "/metrics-alternate")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
