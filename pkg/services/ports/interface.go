package ports

import (
	"context"
	"errors"
)

var (
	// ErrPortNotFound is returned by PortService methods when port is not found.
	ErrPortNotFound = errors.New("port not found")
	// ErrPortAlreadyExist is returned by PortService methods when port is not found.
	ErrPortAlreadyExist = errors.New("port already exist")
)

// Port describes port specific information.
type Port struct {
	// City is a city of a port.
	City string
	// Coordinates is a coordinates of a port.
	Coordinates []float64
	// Country is a country of a port.
	Country string
	// Name is a name of a port.
	Name string
	// Province is a province of a port.
	Province string
}

// PortService is a port service interface.
type PortService interface {
	// Create creates a new port entry.
	Create(ctx context.Context, ID string, port Port) error
	// Get returns port for a given port's ID.
	Get(_ context.Context, ID string) (Port, error)
	// Update updates an existing port.
	Update(ctx context.Context, ID string, port Port) error
}
