package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ExampleGetPort is an example how to get existing port.
func ExampleGetPort() {
	portID := "ExampleCreatePort_" + randString(6)
	endpoint := fmt.Sprintf("http://localhost:%s/api/v1/ports/%s", apiPort, portID)
	http.Get(endpoint) // nolint: noctx
}

func TestGetPort(t *testing.T) {
	portID := "TestGetPort_" + randString(6)

	passed := t.Run("port does not exit", func(t *testing.T) {
		resp, err := http.Get(portsService + "/" + portID) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	require.True(t, passed)
}
