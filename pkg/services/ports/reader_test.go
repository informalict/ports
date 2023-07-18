package ports

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadPorts(t *testing.T) {

	t.Run("invalid root object", func(t *testing.T) {
		channel := make(chan PortWithID)
		data := `invalid`
		reader := bytes.NewReader([]byte(data))

		err := ReadPorts(context.Background(), reader, channel)
		require.Error(t, err)
	})

	t.Run("empty data", func(t *testing.T) {
		channel := make(chan PortWithID)
		data := `{}`
		reader := bytes.NewReader([]byte(data))

		err := ReadPorts(context.Background(), reader, channel)
		require.NoError(t, err)
	})

	t.Run("get port data sub object", func(t *testing.T) {
		channel := make(chan PortWithID)
		data := `{ "portID": { "name": "portName", "coordinates": [1.0, 2.0] } }`
		reader := bytes.NewReader([]byte(data))

		var port PortWithID
		go func() {
			port = <-channel
		}()
		err := ReadPorts(context.Background(), reader, channel)
		require.NoError(t, err)
		expectedPort := PortWithID{
			Port: Port{
				Coordinates: []float64{1, 2},
				Name:        "portName",
			},
			ID: "portID",
		}
		require.Equal(t, expectedPort, port)
	})

	// TODO test for context interruption.
}
