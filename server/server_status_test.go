//+build integration

package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeverStatusEndpoint(t *testing.T) {
	BeforeEach(t)

	s, _, _ := sut(t)

	defer s.Shutdown(context.Background())

	resp, err := http.Get(HTTPAddress + "/status")
	assert.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	expected := `{"status": "ok", "info": {"build":"af09", "build_time":"1988-01-21", "name":"go-serve", "version":"v1.0.0"}}`
	assert.JSONEq(t, expected, string(data))
}
