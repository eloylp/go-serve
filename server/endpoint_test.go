package server //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_endpointMapper_Map(t *testing.T) {
	em := newEndpointMapper()
	em.Declare("/static", "/static")
	result := em.Map("/static/sub/file.txt")
	require.Equal(t, "/static", result)
}
