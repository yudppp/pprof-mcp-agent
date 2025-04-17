package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	pprofmcpagent "github.com/yudppp/pprof-mcp-agent"
)

func main() {
	// Enable block profiling
	runtime.SetBlockProfileRate(1)

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

// heavyProcess simulates various performance issues for profiling demonstration
func heavyProcess(ctx context.Context) {
	// Create a memory leak with a global variable
	var leakySlice [][]int

	// Create multiple goroutines that will never be cleaned up
	for i := 0; i < 100; i++ {
		go func() {
			select {} // Goroutine leak
		}()
	}

	// Channel to demonstrate blocking operations
	ch := make(chan int)

	// Goroutine for blocking operations
	go func() {
		for {
			ch <- 1 // Will block due to unbuffered channel
			time.Sleep(100 * time.Millisecond)
		}
	}()

	mutex := &sync.Mutex{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			// CPU-intensive operations
			go func() {
				mutex.Lock()
				defer mutex.Unlock()

				// Heavy CPU computation
				result := 0
				for i := 0; i < 1000000; i++ {
					result += i * i
				}
			}()

			// Memory-intensive operations
			data := make([]int, 1000000)
			for i := 0; i < len(data); i++ {
				data[i] = i
			}

			// Memory leak simulation
			leakySlice = append(leakySlice, data)

			// Allocate and deallocate to create memory churn
			for i := 0; i < 1000; i++ {
				temp := make([]byte, 1024*1024)
				_ = temp
			}

			// Thread creation
			go func() {
				runtime.LockOSThread()
				select {} // Never releases the OS thread
			}()

			time.Sleep(100 * time.Millisecond)
		}
	}
}
