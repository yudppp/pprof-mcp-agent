package main

import (
	"context"
	"fmt"
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

	go heavyProcess()

	fmt.Println("MCP server listening on :1239")
	err := pprofmcpagent.ServeSSE(ctx, ":1239")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}

// heavyProcess simulates a CPU-intensive task for profiling demonstration
func heavyProcess() {
	for {
		var data []int
		for i := 0; i < 1000000; i++ {
			data = append(data, i)
		}
		_ = data
		time.Sleep(1 * time.Second)
	}
}
