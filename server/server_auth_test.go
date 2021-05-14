//+build integration

package server_test

import (
	"bytes"
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

var testUserCredentials = map[string]string{
	"user": "$2y$10$mAx10mlJ/UNbQJCgPp2oLe9n9jViYl9vlT0cYI3Nfop3P3bU1PDay", // Unencrypted value: user:password
}

func TestReadAuthorizedUserIsAccepted(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithReadAuthorizations(testUserCredentials),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())

	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodGet, HTTPAddress, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("user", "password")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadNonAuthorizedUserIsRefused(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithReadAuthorizations(testUserCredentials),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())

	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodGet, HTTPAddress, nil)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestReadBadlyAuthorizedUserIsRefused(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithReadAuthorizations(testUserCredentials),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())

	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodGet, HTTPAddress, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("user", "bad-password")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWriteAuthorizedUserIsAccepted(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithWriteAuthorizations(testUserCredentials),
			config.WithUploadEndpoint("/upload"),
			config.WithDocRoot(t.TempDir()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", bytes.NewReader(sampleTARGZContent))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.SetBasicAuth("user", "password")
	respAuth, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer respAuth.Body.Close()
	assert.Equal(t, http.StatusOK, respAuth.StatusCode)
}

func TestWriteNonAuthorizedUserIsRefused(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithWriteAuthorizations(testUserCredentials),
			config.WithUploadEndpoint("/upload"),
			config.WithDocRoot(t.TempDir()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", bytes.NewReader(sampleTARGZContent))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWriteBadlyAuthorizedUserIsRefused(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
			config.WithWriteAuthorizations(testUserCredentials),
			config.WithUploadEndpoint("/upload"),
			config.WithDocRoot(t.TempDir()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", bytes.NewReader(sampleTARGZContent))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.SetBasicAuth("user", "bad-password")
	respAuth, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer respAuth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, respAuth.StatusCode)
}
