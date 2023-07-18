package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/informalict/ports/api/v1"
	"github.com/informalict/ports/pkg/services/ports/memory"
)

var apiPorts = "/api/v1/ports"

func getEndpoint(server *httptest.Server, portID string) string {
	return server.URL + apiPorts + "/" + portID
}

// TestGetPort tests for getting port.
func TestGetPort(t *testing.T) {
	stub := memory.NewPortMemory()
	router := NewPortRouter(stub)
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	portID := "test"

	passed := t.Run("port does not exit", func(t *testing.T) {
		resp, err := client.Get(getEndpoint(server, portID)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	require.True(t, passed)

	passed = t.Run("empty port ID", func(t *testing.T) {
		resp, err := client.Get(getEndpoint(server, "")) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	require.True(t, passed)

	passed = t.Run("create and get port", func(t *testing.T) {
		validPort := api.Port{
			Name:        "name",
			City:        "city",
			Country:     "country",
			Coordinates: []float64{1.0, 1.0},
		}
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		resp, err = client.Get(getEndpoint(server, portID)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		svcPort, err := ParseRequestPort(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, validPort, svcPort)
	})
	require.True(t, passed)
}

// TestCreatePort tests for port creation.
func TestCreatePort(t *testing.T) { // nolint: funlen
	stub := memory.NewPortMemory()
	router := NewPortRouter(stub)
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	portID := "test"

	passed := t.Run("failed to parse input data", func(t *testing.T) {
		invalidInput := "test"
		b, err := json.Marshal(&invalidInput)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errorMsg, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		expectedError := errors.New(strings.TrimSuffix(string(errorMsg), "\n"))
		require.EqualError(t, expectedError, "failed to parse input port's data")
	})
	require.True(t, passed)

	passed = t.Run("port's name can not be empty", func(t *testing.T) {
		validPort := api.Port{
			Name: "",
		}
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errorMsg, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		expectedError := errors.New(strings.TrimSuffix(string(errorMsg), "\n"))
		require.EqualError(t, expectedError, "port's name can not be empty")
	})
	require.True(t, passed)

	validPort := api.Port{
		Name:        "name",
		City:        "city",
		Country:     "country",
		Coordinates: []float64{1.0, 1.0},
	}

	passed = t.Run("create new port", func(t *testing.T) {
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})
	require.True(t, passed)

	passed = t.Run("can not create existing port", func(t *testing.T) {
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusConflict, resp.StatusCode)
	})
	require.True(t, passed)
}

// TestUpdatePort tests for port update.
func TestUpdatePort(t *testing.T) { // nolint: funlen
	stub := memory.NewPortMemory()
	router := NewPortRouter(stub)
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()
	portID := "test"

	passed := t.Run("failed to parse input data", func(t *testing.T) {
		invalidInput := "test"
		b, err := json.Marshal(&invalidInput)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, getEndpoint(server, portID), bytes.NewBuffer(b)) // nolint: noctx
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		errorMsg, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		expectedError := errors.New(strings.TrimSuffix(string(errorMsg), "\n"))
		require.EqualError(t, expectedError, "failed to parse input port's data")
	})
	require.True(t, passed)

	passed = t.Run("coordinates' can not be empty", func(t *testing.T) {
		validPort := api.Port{
			Name: "test",
		}
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errorMsg, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		expectedError := errors.New(strings.TrimSuffix(string(errorMsg), "\n"))
		require.EqualError(t, expectedError, "port's coordinates can not be empty")
	})
	require.True(t, passed)

	validPort := api.Port{
		Name:        "name",
		City:        "city",
		Country:     "country",
		Coordinates: []float64{1.0, 1.0},
	}

	passed = t.Run("create new port", func(t *testing.T) {
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)
		resp, err := client.Post(getEndpoint(server, portID), "application/json", bytes.NewReader(b)) // nolint: noctx
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})
	require.True(t, passed)

	passed = t.Run("can not create existing port", func(t *testing.T) {
		b, err := json.Marshal(&validPort)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, getEndpoint(server, portID), bytes.NewBuffer(b)) // nolint: noctx
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	require.True(t, passed)
}
