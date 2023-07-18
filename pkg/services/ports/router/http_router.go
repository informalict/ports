package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	api "github.com/informalict/ports/api/v1"
	"github.com/informalict/ports/pkg/services/ports"
)

const (
	apiV1Prefix = "/api/v1/"
)

// portRouter describes HTTP router for a port service.
type portRouter struct {
	svc ports.PortService
}

// NewPortRouter returns a new port's router for a given port service.
func NewPortRouter(svc ports.PortService) http.Handler {
	pr := &portRouter{
		svc,
	}

	router := httprouter.New()
	router.GET(apiV1Prefix+"ports/:id", pr.GetPort)
	router.POST(apiV1Prefix+"ports/:id", pr.CreatePort)
	router.PUT(apiV1Prefix+"ports/:id", pr.UpdatePort)

	return router
}

// UpdatePort updates a port in a storage.
func (pr *portRouter) UpdatePort(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		http.Error(w, "id of a port must be provided", http.StatusBadRequest)
		return
	}

	apiPort, err := ParseRequestPort(r.Body)
	if err != nil {
		// It should be error log level.
		log.Println(fmt.Sprintf("failed to parse input port's data: %s\n", err))
		http.Error(w, "failed to parse input port's data", http.StatusInternalServerError)
		return
	}

	if err := apiPort.Validate(); err != nil {
		// It should be error log level.
		log.Println(fmt.Sprintf("failed to validate input port's data: %s\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port := convertFromAPIPort(apiPort)
	if err := pr.svc.Update(r.Context(), id, port); err != nil {
		if errors.Is(err, ports.ErrPortNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			// It should be error log level.
			log.Println(fmt.Sprintf("failed to update a port: %s\n", err))
			http.Error(w, "failed to update a port", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// CreatePort creates a new port in a storage.
func (pr *portRouter) CreatePort(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		http.Error(w, "id of a port must be provided", http.StatusBadRequest)
		return
	}

	apiPort, err := ParseRequestPort(r.Body)
	if err != nil {
		// It should be error log level.
		log.Println(fmt.Sprintf("failed to parse input port's data: %s\n", err))
		http.Error(w, "failed to parse input port's data", http.StatusBadRequest)
		return
	}

	if err := apiPort.Validate(); err != nil {
		// It should be error log level.
		log.Println(fmt.Sprintf("failed to validate input port's data: %s\n", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port := convertFromAPIPort(apiPort)
	if err := pr.svc.Create(r.Context(), id, port); err != nil {
		if errors.Is(err, ports.ErrPortAlreadyExist) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			// It should be error log level.
			log.Println(fmt.Sprintf("failed to create a new port: %s\n", err))
			http.Error(w, "failed to create a new port", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return
}

// GetPort is an HTTP handler which fetches port from a port's service.
func (pr *portRouter) GetPort(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		http.Error(w, "id of a port must be provided", http.StatusBadRequest)
		return
	}

	port, err := pr.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ports.ErrPortNotFound) {
			// No logs or it can be debug log level.
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// It should be error log level.
		log.Println(fmt.Sprintf("failed to get port: %s\n", err))
		http.Error(w, "failed to get port", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(ConvertToAPIPort(port))
	if err != nil {
		// It should be error log level.
		log.Println(fmt.Sprintf("failed to marhal port: %s", err))
		http.Error(w, "failed to serialize a port", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Println(err)
	}
}

// ConvertToAPIPort converts internal port structure to client api structure.
func ConvertToAPIPort(port ports.Port) api.Port {
	return api.Port{
		City:        port.City,
		Coordinates: port.Coordinates,
		Country:     port.Country,
		Name:        port.Name,
		Province:    port.Province,
	}
}

// convertFromAPIPort converts client API port into internal API port.
func convertFromAPIPort(port api.Port) ports.Port {
	return ports.Port{
		City:        port.City,
		Coordinates: port.Coordinates,
		Country:     port.Country,
		Name:        port.Name,
		Province:    port.Province,
	}
}

// ParseRequestPort validates port from a caller.
func ParseRequestPort(r io.Reader) (api.Port, error) {
	var apiPort api.Port

	body, err := io.ReadAll(r)
	if err != nil {
		return apiPort, err
	}

	if err := json.Unmarshal(body, &apiPort); err != nil {
		return apiPort, err
	}

	return apiPort, nil
}
