package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve"
	"github.com/eloylp/go-serve/config"
)

const (
	SutListenAddress = "localhost:9090"
	SutHTTPAddress   = "http://" + SutListenAddress
	TestDocRoot      = "tests/root"
	TuxTestFileMD5   = "a0e6e27f7e31fd0bd549ea936033bf28"
)

func TestServingContent(t *testing.T) {
	cfg := config.ForOptions(
		config.WithListenAddr(SutListenAddress),
		config.WithDocRoot(TestDocRoot),
		config.WithDocRootPrefix("/"),
	)
	s := server.New(cfg)
	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	data := BodyFrom(t, SutHTTPAddress+"/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
}
