package server_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

const (
	ListenAddress  = "localhost:9090"
	HTTPAddress    = "http://" + ListenAddress
	DocRoot        = "../tests/root"
	TuxTestFileMD5 = "a0e6e27f7e31fd0bd549ea936033bf28"
)

func init() {
	server.Name = "programName"
	server.Version = "v1.0.0"
	server.Build = "af09"
	server.BuildTime = "1988-01-21"
}

func TestServingContent(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	cfg := config.ForOptions(
		config.WithListenAddr(ListenAddress),
		config.WithDocRoot(DocRoot),
		config.WithDocRootPrefix("/"),
		config.WithLoggerOutput(logBuff),
	)
	s, err := server.New(cfg)
	assert.NoError(t, err)
	go s.ListenAndServe()
	data := BodyFrom(t, HTTPAddress+"/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}

func AssertNoProblemsInLogs(t *testing.T, logs string) {
	assert.NotContains(t, logs, "level=warning")
	assert.NotContains(t, logs, "level=error")
}

func AssertStartupLogs(t *testing.T, logs string) {
	assert.Contains(t, logs, "programName v1.0.0 af09 1988-01-21")
	assert.Contains(t, logs, fmt.Sprintf("Starting to serve"))
}

func AssertShutdownLogs(t *testing.T, logs string) {
	assert.Contains(t, logs, "started gracefully shutdown of server ...")
	assert.Contains(t, logs, "server is now shutdown !")
}