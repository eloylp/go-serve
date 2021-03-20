package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eloylp/kit/test"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestSeverIdentity(t *testing.T) {
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

	resp, err := http.Get(HTTPAddress)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "go-serve v1.0.0", resp.Header.Get("server"))
}
