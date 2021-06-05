//+build integration

package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeverIdentity(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

	resp, err := http.Get(HTTPAddressStatic)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "go-serve v1.0.0", resp.Header.Get("server"))
}
