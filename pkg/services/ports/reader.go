package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// PortWithID extends Port structure with ID field.
type PortWithID struct {
	Port
	ID string
}

// ReadPorts reads port's data from a given reader and sends it to the channel.
// When iterating JSON data is faster than sending to the channel then it will hang until caller receives data.
// A caller can provide buffered channel, so it can control how many messages are in a memory.
// A reader must provide data in valid format `{ "portID1": {}, "portID2": {}, ... }`.
func ReadPorts(ctx context.Context, reader io.Reader, channel chan<- PortWithID) error {
	decoder := json.NewDecoder(reader)
	// Go to the first entry in a map.
	if _, err := decoder.Token(); err != nil {
		return err
	}

	var portID string
	var validID bool
	for decoder.More() {
		if ctx.Err() != nil {
			// A caller triggered cancellation.
			return ctx.Err()
		}

		// Get port ID for a next entry.
		if v, err := decoder.Token(); err != nil {
			return err
		} else {
			if portID, validID = v.(string); !validID {
				return fmt.Errorf("string type is expected for port keys, but got \"%T\"", v)
			}
		}

		var port PortWithID
		if err := decoder.Decode(&port.Port); err == nil {
			port.ID = portID
			channel <- port
		}
	}

	return nil
}
