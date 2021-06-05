//+build integration

package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverIdentity(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

	resp, err := http.Get(HTTPAddressStatic)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "go-serve v1.0.0", resp.Header.Get("server"))
}
