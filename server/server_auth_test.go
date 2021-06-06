//+build integration

package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.eloylp.dev/go-serve/config"
)

var testUserCredentials = map[string]string{
	"user": "$2y$10$mAx10mlJ/UNbQJCgPp2oLe9n9jViYl9vlT0cYI3Nfop3P3bU1PDay", // Unencrypted value: user:password
}

func TestReadAuthorizedUserIsAccepted(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithReadAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodGet, HTTPAddressStatic, nil)
	require.NoError(t, err)
	req.SetBasicAuth("user", "password")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadNonAuthorizedUserIsRefused(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithReadAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodGet, HTTPAddressStatic, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestReadBadlyAuthorizedUserIsRefused(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithReadAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodGet, HTTPAddressStatic, nil)
	require.NoError(t, err)
	req.SetBasicAuth("user", "bad-password")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWriteAuthorizedUserIsAccepted(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithWriteAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodPost, HTTPAddressUpload, sampleTARGZContentReader())
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.SetBasicAuth("user", "password")
	respAuth, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer respAuth.Body.Close()
	assert.Equal(t, http.StatusOK, respAuth.StatusCode)
}

func TestWriteNonAuthorizedUserIsRefused(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithWriteAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodPost, HTTPAddressUpload, sampleTARGZContentReader())
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWriteBadlyAuthorizedUserIsRefused(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t, config.WithWriteAuthorizations(testUserCredentials))

	defer s.Shutdown(context.Background())

	req, err := http.NewRequest(http.MethodPost, HTTPAddressUpload, sampleTARGZContentReader())
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.SetBasicAuth("user", "bad-password")
	respAuth, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer respAuth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, respAuth.StatusCode)
}
