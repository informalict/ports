package test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/informalict/ports/pkg/services/ports"
	"github.com/informalict/ports/pkg/services/ports/router"
)

// TestCheckInitFile check whether all entries from file exist in the service.
// This test is idempotent and can be launched many times.
func TestCheckInitFile(t *testing.T) {
	channel := make(chan ports.PortWithID)

	file, err := os.Open(testFile)
	require.NoError(t, err)

	go func() {
		for {
			select {
			case filePort, ok := <-channel:
				if !ok {
					// No data and channel is closed, so all data is fetched.
					return
				}

				resp, err := http.Get(portsService + "/" + filePort.ID) // nolint: noctx
				require.NoError(t, err)

				svcPort, err := router.ParseRequestPort(resp.Body)
				resp.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, router.ConvertToAPIPort(filePort.Port), svcPort)
			}
		}
	}()

	err = ports.ReadPorts(context.Background(), file, channel)
	require.NoError(t, err)
}
