//+build integration

package server_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestServingContentDefaults(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(DocRoot),
			config.WithLoggerOutput(logBuff),
		),
	)
	assert.NoError(t, err)
	go s.ListenAndServe()

	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	data := BodyFrom(t, HTTPAddressStatic+"/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}

func TestServingContentAlternatePath(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(DocRoot),
			config.WithDocRootPrefix("/alternate"),
			config.WithLoggerOutput(logBuff),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	data := BodyFrom(t, HTTPAddress+"/alternate/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}
