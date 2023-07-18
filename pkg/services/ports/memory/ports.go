package memory

import (
	"context"
	"sync"

	"github.com/informalict/ports/pkg/services/ports"
)

// NewPortMemory creates port's memory storage.
func NewPortMemory() *portMemory {
	return &portMemory{
		ports: make(map[string]ports.Port),
	}
}

type portMemory struct {
	// mutex locks this structure when CRUD actions are performed.
	mutex sync.RWMutex
	// ports stores ports.
	ports map[string]ports.Port
}

// Create creates a port in memory with a given port ID.
func (p *portMemory) Create(_ context.Context, ID string, port ports.Port) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.ports[ID]; ok {
		return ports.ErrPortAlreadyExist
	}
	p.ports[ID] = port

	return nil
}

// Get returns port for a given port's ID.
func (p *portMemory) Get(_ context.Context, ID string) (ports.Port, error) {
	if port, ok := p.ports[ID]; ok {
		return port, nil
	}

	return ports.Port{}, ports.ErrPortNotFound
}

// Update updates an existing port.
// When port does not exist then error is returned.
func (p *portMemory) Update(_ context.Context, ID string, port ports.Port) error {
	if _, ok := p.ports[ID]; !ok {
		return ports.ErrPortNotFound
	}

	p.ports[ID] = port

	return nil
}
