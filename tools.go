package pprofmcpagent

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Package pprofmcpagent provides a Mission Control Protocol (MCP) agent for Go's runtime profiling.
// It enables real-time collection and analysis of various Go runtime profiles through MCP.
// The agent supports multiple view modes (flat, cumulative, and graph) and various profile types
// for comprehensive performance analysis.

// NewPprofServer creates a new MCP server with all pprof tools registered.
// It initializes a server instance with all available profiling tools:
// - Heap profiling (memory allocations)
// - Goroutine profiling (stack traces)
// - Thread creation profiling (OS threads)
// - Block profiling (synchronization)
// - Allocation profiling (memory usage)
// - CPU profiling (execution time)
//
// Each profile can be viewed in three modes:
// - Flat: direct values for each function
// - Cumulative: including child function costs
// - Graph: showing call relationships
func NewPprofServer() *server.MCPServer {
	s := server.NewMCPServer(
		"pprof server",
		"1.0.0",
	)

	// Add tools
	s.AddTool(NewHeapTool(), HeapHandler)
	s.AddTool(NewGoroutineTool(), GoroutineHandler)
	s.AddTool(NewThreadCreateTool(), ThreadCreateHandler)
	s.AddTool(NewBlockTool(), BlockHandler)
	s.AddTool(NewAllocsTool(), AllocsHandler)
	s.AddTool(NewCPUTool(), CPUHandler)

	return s
}

// newProfileTool creates a new MCP tool with common profile configuration options.
// It sets up standard parameters like the result limit and view mode, along with any
// additional tool-specific options provided.
//
// Parameters:
//   - name: The name of the profiling tool
//   - description: A description of what the tool does
//   - extraOpts: Additional tool-specific options
//
// Configuration options:
//   - limit: Number of top locations to show (default: 100, min: 100, max: 10000)
//   - view: Profile view mode (flat, cum, graph)
func newProfileTool(name, description string, extraOpts ...mcp.ToolOption) mcp.Tool {
	opts := []mcp.ToolOption{
		mcp.WithDescription(description),
		mcp.WithNumber(
			"limit",
			mcp.Description("Maximum number of locations to show in results"),
			mcp.DefaultNumber(100),
			mcp.Min(100),
			mcp.Max(10000),
		),
		mcp.WithString(
			"view",
			mcp.Description("View mode for profile data (flat: direct values, cum: cumulative values including children, graph: call graph)"),
			mcp.DefaultString(string(ViewModeFlat)),
			mcp.Enum(
				string(ViewModeFlat),
				string(ViewModeCum),
				string(ViewModeGraph),
			),
		),
	}
	opts = append(opts, extraOpts...)
	return mcp.NewTool(name, opts...)
}

// NewHeapTool creates a new MCP tool for heap profiling.
// This tool provides insights into memory usage patterns and potential memory leaks.
// It shows current memory allocations by location and helps identify inefficient memory usage.
func NewHeapTool() mcp.Tool {
	return newProfileTool("heap-profile", "Output heap memory profile data")
}

// NewGoroutineTool creates a new MCP tool for goroutine profiling.
// This tool helps identify goroutine leaks and analyze concurrency patterns
// by providing detailed stack traces of currently running goroutines.
func NewGoroutineTool() mcp.Tool {
	return newProfileTool("goroutine-profile", "Output goroutine stack traces")
}

// NewThreadCreateTool creates a new MCP tool for thread creation profiling.
// This tool helps track OS thread creation patterns and potential thread leaks,
// useful for analyzing thread pool behavior and system resource usage.
func NewThreadCreateTool() mcp.Tool {
	return newProfileTool("threadcreate-profile", "Output thread creation profile data")
}

// NewBlockTool creates a new MCP tool for block profiling.
// This tool helps identify synchronization bottlenecks and deadlock risks
// by showing where goroutines block on synchronization primitives.
func NewBlockTool() mcp.Tool {
	return newProfileTool("block-profile", "Output blocking operation profile data")
}

// NewAllocsTool creates a new MCP tool for allocation profiling.
// This tool helps analyze memory allocation patterns and identify memory churn
// by showing both allocated and freed memory statistics.
func NewAllocsTool() mcp.Tool {
	return newProfileTool("allocs-profile", "Output memory allocation sampling data")
}

// NewCPUTool creates a new MCP tool for CPU profiling.
// This tool helps identify CPU-intensive code paths and performance bottlenecks
// by sampling program execution over a specified duration.
func NewCPUTool() mcp.Tool {
	return newProfileTool("cpu-profile", "Output CPU profile data",
		mcp.WithNumber(
			"duration",
			mcp.Description("Duration of CPU profiling in seconds"),
			mcp.DefaultNumber(10),
		),
	)
}
