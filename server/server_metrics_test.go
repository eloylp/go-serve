//+build integration

package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
)

func TestMetricsAreServedByDefault(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

	resp, err := http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMetricsAreObserving(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

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
	assert.Contains(t, metrics, `http_request_duration_seconds_count{code="200",endpoint="/static",method="GET"} 1`)
}

func TestMetricsCanBeDisabled(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithMetricsEnabled(false))

	defer s.Shutdown(context.Background())

	AssertDefaultMetricsPathConfigured(t)
}

func AssertDefaultMetricsPathConfigured(t *testing.T) {
	resp, err := http.Get(HTTPAddress + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMetricsCanBeServedOnAlternativePort(t *testing.T) {
	BeforeEach(t)

	s, loggerOutput, _ := sut(t, config.WithMetricsAlternativeListenAddr("0.0.0.0:9091"))

	defer s.Shutdown(context.Background())

	test.WaitTCPService(t, "localhost:9091", time.Millisecond, time.Second)

	AssertDefaultMetricsPathConfigured(t)

	resp, err := http.Get("http://localhost:9091/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, loggerOutput.String(), "starting to serve metrics at 0.0.0.0:9091")
}

func TestMetricsRequestDurationBucketsConfig(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

	// Make a request.
	resp, err := http.Get(HTTPAddressStatic)
	require.NoError(t, err)
	defer resp.Body.Close()

	resp, err = http.Get(HTTPAddress + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	metrics := string(data)

	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",endpoint="/static",method="GET",le="0.1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",endpoint="/static",method="GET",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",endpoint="/static",method="GET",le="1"} 1`)
}

func TestMetricsResponseSizeBucketsConfig(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithMetricsSizeBuckets([]float64{1, 4, 6}))

	defer s.Shutdown(context.Background())

	// Make a request.
	resp, err := http.Get(HTTPAddressStatic)
	require.NoError(t, err)
	defer resp.Body.Close()

	resp, err = http.Get(HTTPAddress + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	metrics := string(data)

	require.Contains(t, metrics, `http_response_size_bucket{code="200",endpoint="/static",method="GET",le="1"} 0`)
	require.Contains(t, metrics, `http_response_size_bucket{code="200",endpoint="/static",method="GET",le="4"} 0`)
	require.Contains(t, metrics, `http_response_size_bucket{code="200",endpoint="/static",method="GET",le="6"} 0`)
}

func TestMetricsUploadSizeBucketsConfig(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithMetricsSizeBuckets([]float64{100_000, 600_000, 1_000_000}))

	defer s.Shutdown(context.Background())

	tarGZFile, err := os.Open(DocRootTARGZ)
	require.NoError(t, err)
	defer tarGZFile.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", tarGZFile)
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.Header.Add("GoServe-Deploy-Path", "/deploy")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	resp, err = http.Get(HTTPAddress + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	metrics := string(data)

	require.Contains(t, metrics, `goserve_upload_size_bucket{le="100000"} 0`)
	require.Contains(t, metrics, `goserve_upload_size_bucket{le="600000"} 1`)
	require.Contains(t, metrics, `goserve_upload_size_bucket{le="1e+06"} 1`)
}

func TestMetricsCanBeServedAlternativePath(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithMetricsPath("/metrics-alternate"))

	defer s.Shutdown(context.Background())

	resp, err := http.Get(HTTPAddress + "/metrics-alternate")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
