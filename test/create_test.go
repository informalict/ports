package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/informalict/ports/api/v1"
	"github.com/informalict/ports/pkg/services/ports/router"
)

// TestCreatePort creates a new port.
// This test may not be idempotent, because it can not be cleaned up.
func TestCreatePort(t *testing.T) {
	portID := "TestCreatePort_" + randString(6)

	passed := t.Run("name can not be empty", func(t *testing.T) {
		invalidPort := api.Port{}
		b, err := json.Marshal(&invalidPort)
		require.NoError(t, err)

		resp, err := http.Post(portsService+"/"+portID, "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	require.True(t, passed)

	passed = t.Run("create a port", func(t *testing.T) {
		validPort := api.Port{
			Name:        "name",
			City:        "city",
			Country:     "country",
			Coordinates: []float64{1.0, 1.0},
		}
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)

		resp, err := http.Post(portsService+"/"+portID, "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		resp, err = http.Get(portsService + "/" + portID) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		svcPort, err := router.ParseRequestPort(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, validPort, svcPort)
	})
	require.True(t, passed)
}
