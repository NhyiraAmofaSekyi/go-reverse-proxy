package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/configs"
)

var (
	server *http.Server
)

func ConfigRun(cfg *configs.Configuration) error {
	// Load configurations from the config file.

	// Create a new router.
	mux := http.NewServeMux()

	// Register the health check endpoint.
	mux.HandleFunc("/ping", ping)

	// Register configured routes.
	for _, resource := range cfg.Resources {
		destURL, _ := url.Parse(resource.Destination_URL)
		proxy := NewProxy(destURL)
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(proxy, destURL, resource.Endpoint))
	}

	// Initialize the HTTP server.
	server = &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Listen_port,
		Handler: mux,
	}

	fmt.Printf("Server configured to listen on %s\n", server.Addr)

	// Create a channel to listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	// Register the given channel to receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine.
	go func() {
		fmt.Println("Server goroutine starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Block until a signal is received.
	<-quit
	fmt.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown process.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
	return nil
}
