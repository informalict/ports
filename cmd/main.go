package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/informalict/ports/pkg/services/ports"
	"github.com/informalict/ports/pkg/services/ports/memory"
	"github.com/informalict/ports/pkg/services/ports/router"
)

// The below const params should be provided from processes' arguments.
const (
	// portsInMemory describes how many ports can be kept in memory simultaneously.
	portsInMemory = 10
	// initialInputFileName is a file with a initial input data.
	initialInputFileName = "./assets/ports.json"
	// addressApp for the server to listen on.
	addressApp = ":8080"
)

func main() {
	portService := memory.NewPortMemory()
	ctx := createSignalContext()

	if err := readInitFile(ctx, portService); err != nil {
		log.Fatal(err.Error())
	}

	// Start HTTP server.
	srv := &http.Server{
		Addr:              addressApp,
		Handler:           router.NewPortRouter(portService),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err)
		}
	}()

	log.Println("server is ready")
	// Wait until process gets signal SIGTERM.
	select {
	case <-ctx.Done():
		// Perform actions for graceful shutdown.
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(timeoutCtx); err != nil {
			log.Fatalf("failed to shutdown server: %s", err)
		}
	}

	return
}

// createSignalContext creates context which is canceled when SIGTERM occurs.
func createSignalContext() context.Context {
	sigChannel := make(chan os.Signal, 1)
	// SIGKILL can not be handled.
	signal.Notify(sigChannel, syscall.SIGTERM)

	ctxCancel, cancel := context.WithCancel(context.Background())

	go func() {
		// Wait until signal occurs.
		<-sigChannel
		// Close the context which is used by the caller.
		cancel()
		close(sigChannel)
	}()

	return ctxCancel
}

// readInitFile reads data from fixed input file and populate them into port's service.
func readInitFile(ctx context.Context, svc ports.PortService) error {
	channel := make(chan ports.PortWithID, portsInMemory)

	file, err := os.Open(initialInputFileName)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				// Perform actions for graceful shutdown when initial file is not fully loaded.
				return
			case port, ok := <-channel:
				if !ok {
					// No data and channel is closed, so all data is fetched.
					return
				}

				if err := svc.Create(ctx, port.ID, port.Port); err != nil {
					if ctx.Err() != nil {
						continue
					}

					// In real life scenario this case should be handled, and it should not end this process.
					log.Fatalf("failed to add port \"%s\"", port.ID)
				}
			}
		}
	}()

	if err := ports.ReadPorts(ctx, file, channel); err != nil {
		log.Fatal(err)
	}

	// All data is fetched from JSON file, so channel can be closed.
	close(channel)

	return nil
}
