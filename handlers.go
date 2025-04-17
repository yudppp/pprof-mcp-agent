package pprofmcpagent

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime/pprof"
	"time"

	"github.com/google/pprof/profile"
	"github.com/mark3labs/mcp-go/mcp"
)

// Package pprofmcpagent provides a Model Context Protocol (MCP) agent for Go's runtime profiling.
// It enables real-time collection and analysis of various Go runtime profiles through MCP.

// Profile type constants
const (
	ProfileTypeHeap         = "heap"
	ProfileTypeGoroutine    = "goroutine"
	ProfileTypeThreadCreate = "threadcreate"
	ProfileTypeBlock        = "block"
	ProfileTypeAllocs       = "allocs"
	ProfileTypeCPU          = "cpu"
)

// Profile error definitions
type ProfileError struct {
	ProfileType string
	Err         error
}

func (e *ProfileError) Error() string {
	return fmt.Sprintf("profile error (%s): %v", e.ProfileType, e.Err)
}

// handleProfile is a common function that processes various types of runtime profiles.
// It handles profile data collection, parsing, and formatting the results.
func handleProfile(profileName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	prof := pprof.Lookup(profileName)
	if prof == nil {
		return nil, &ProfileError{
			ProfileType: profileName,
			Err:         fmt.Errorf("profile not found"),
		}
	}

	// Get limit from request parameters
	limit := 100
	if limitParam, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(limitParam)
	}

	// Get view mode from request parameters
	viewMode := ViewModeFlat
	if viewParam, ok := request.Params.Arguments["view"].(string); ok {
		viewMode = ViewMode(viewParam)
	}

	// Write profile data to buffer
	var buf bytes.Buffer
	if err := prof.WriteTo(&buf, 0); err != nil {
		return nil, &ProfileError{
			ProfileType: profileName,
			Err:         fmt.Errorf("failed to write profile: %w", err),
		}
	}

	// Parse the profile
	p, err := profile.Parse(&buf)
	if err != nil {
		return nil, &ProfileError{
			ProfileType: profileName,
			Err:         fmt.Errorf("failed to parse profile: %w", err),
		}
	}

	result := getTopSamples(p, limit, viewMode, profileName)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

// HeapHandler processes heap profile requests.
// It provides aggregated memory allocation statistics from the heap,
// showing memory usage by location in the code.
// The results include both in-use and allocated memory statistics.
func HeapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleProfile(ProfileTypeHeap, request)
}

// GoroutineHandler processes goroutine profile requests.
// It provides aggregated statistics about currently running goroutines,
// including their current state (running, waiting, blocked) and stack traces.
func GoroutineHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleProfile(ProfileTypeGoroutine, request)
}

// ThreadCreateHandler processes thread creation profile requests.
// It provides aggregated statistics about OS thread creation,
// showing locations where new OS threads are created and their frequency.
func ThreadCreateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleProfile(ProfileTypeThreadCreate, request)
}

// BlockHandler processes block profile requests.
// It provides aggregated statistics about goroutine blocking operations,
// showing locations where goroutines block on synchronization primitives
// (mutexes, channels, etc.) and the duration of blocking.
func BlockHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleProfile(ProfileTypeBlock, request)
}

// AllocsHandler processes allocation profile requests.
// It provides aggregated memory allocation statistics,
// showing locations where memory allocations occur and their frequency.
// This includes both allocated and freed memory.
func AllocsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleProfile(ProfileTypeAllocs, request)
}

// CPUHandler processes CPU profile requests.
// It collects and provides aggregated CPU usage statistics over a specified duration,
// showing where the program spends its CPU time.
// The duration can be configured through the request parameters (default: 10 seconds).
func CPUHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	duration, ok := request.Params.Arguments["duration"].(float64)
	if !ok {
		duration = 10
	}

	// Get limit from request parameters
	limit := 100
	if limitParam, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(limitParam)
	}

	var buf bytes.Buffer
	if err := pprof.StartCPUProfile(&buf); err != nil {
		return handleMCPError(err), nil
	}
	time.Sleep(time.Duration(duration) * time.Second)
	pprof.StopCPUProfile()

	// Parse the profile
	p, err := profile.Parse(&buf)
	if err != nil {
		return handleMCPError(err), nil
	}
	result := getTopSamples(p, limit, ViewModeFlat, ProfileTypeCPU)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(result),
		},
	}, nil
}

// handleMCPError creates an error response for MCP tool requests.
func handleMCPError(err error) *mcp.CallToolResult {
	log.Println(err)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(err.Error()),
		},
		IsError: true,
	}
}
