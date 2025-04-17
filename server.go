package pprofmcpagent

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// ServeSSE starts a server that exposes pprof data through Server-Sent Events (SSE).
// The server runs until the provided context is cancelled.
//
// Parameters:
//   - ctx: Context for controlling the server lifecycle
//   - port: Port number to listen on (e.g., ":8080")
//
// Returns:
//   - error: Any error that occurred during server startup or operation
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	if err := ServeSSE(ctx, ":8080"); err != nil {
//	    log.Fatal(err)
//	}
func ServeSSE(ctx context.Context, port string) error {
	s := NewPprofServer()

	// Configure SSE server with enhanced settings
	sses := server.NewSSEServer(s,
		server.WithBaseURL(
			fmt.Sprintf("http://localhost%s", port),
		),
	)

	// Setup graceful shutdown
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		sses.Shutdown(shutdownCtx)
	}()

	if err := sses.Start(port); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
