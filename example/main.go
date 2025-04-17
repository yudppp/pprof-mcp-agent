package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pprofmcpagent "github.com/yudppp/pprof-mcp-agent"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	go heavyProcess(ctx)

	log.Println("MCP server listening on :1239")
	err := pprofmcpagent.ServeSSE(ctx, ":1239")
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		return
	}
}

// heavyProcess simulates a CPU-intensive task for profiling demonstration
func heavyProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var data []int
			for i := 0; i < 1000000; i++ {
				data = append(data, i)
			}
			_ = data
			time.Sleep(1 * time.Second)
		}
	}
}
