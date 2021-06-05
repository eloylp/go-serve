//+build integration

package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
)

func TestServingContentDefaults(t *testing.T) {
	BeforeEach(t)

	s, logBuff, docRoot := sut(t)

	test.Copy(t, DocRoot, docRoot)

	data := BodyFrom(t, HTTPAddressStatic+"/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err := s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()

	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}

func TestServingContentAlternatePath(t *testing.T) {
	BeforeEach(t)

	s, logBuff, docRoot := sut(t, config.WithDocRootPrefix("/alternate"))

	test.Copy(t, DocRoot, docRoot)

	data := BodyFrom(t, HTTPAddress+"/alternate/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err := s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}
