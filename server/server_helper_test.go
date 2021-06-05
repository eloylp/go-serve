//+build integration

package server_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func BeforeEach(_ *testing.T) {
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
}

func BodyFrom(t *testing.T, url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func md5From(data []byte) string {
	md5.New()
	return fmt.Sprintf("%x", md5.Sum(data))
}

// sut will retrieve an started, ready to use, go-serve server instance.
// It will use a standard configuration. Config can be overwritten by passing
// extra config options in the variadic part.
//
// By default it will use a default t.TempDir() as document root and will log
// all logs in debug mode.
//
// It will return the server instance in order to properly defer the shutdown
// process. The buffer where the logs will be written and the current document
// root , in case the test requires some initial state.
//
// More constants and details are defined at source file server_test.go.
func sut(t *testing.T, options ...config.Option) (*server.Server, *bytes.Buffer, string) {
	logs := bytes.NewBuffer(nil)
	docRoot := t.TempDir()
	o := []config.Option{
		config.WithListenAddr(ListenAddress),
		config.WithLoggerOutput(logs),
		config.WithUploadEndpoint("/upload"),
		config.WithDownLoadEndpoint("/download"),
		config.WithDocRoot(docRoot),
		config.WithLoggerLevel(logrus.DebugLevel.String()),
	}
	o = append(o, options...)
	srv, err := server.New(config.ForOptions(o...))
	if err != nil {
		t.Fatal(err)
	}
	go srv.ListenAndServe()
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)
	return srv, logs, docRoot
}
