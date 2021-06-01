//+build integration

package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.eloylp.dev/kit/test"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestSeverStatusEndpoint(t *testing.T) {
	BeforeEach(t)
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

	resp, err := http.Get(HTTPAddress + "/status")
	assert.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	expected := `{"status": "ok", "info": {"build":"af09", "build_time":"1988-01-21", "name":"go-serve", "version":"v1.0.0"}}`
	assert.JSONEq(t, expected, string(data))
}
