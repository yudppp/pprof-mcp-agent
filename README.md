# pprof-mcp-agent

A Go agent that provides runtime profiling data through the Model Context Protocol (MCP). It enables real-time collection and analysis of Go runtime performance data, making it easier to monitor and debug Go applications in production.

## Installation

```bash
go get github.com/yudppp/pprof-mcp-agent
```

## Overview

pprof-mcp-agent integrates Go's built-in profiling capabilities with the MCP protocol, providing a streamlined way to collect and analyze performance data. It supports various profile types, each designed to help diagnose different aspects of your application's performance:

### Supported Profile Types

- **CPU Profile**: Identifies CPU-intensive code paths and performance bottlenecks
  - Configurable sampling duration (default: 10 seconds)
  - Shows where your program spends its execution time
  - Useful for optimizing CPU-bound applications

- **Heap Profile**: Analyzes memory usage patterns
  - Shows current memory allocations by location
  - Helps identify memory leaks and inefficient memory usage
  - Includes both in-use and allocated memory statistics

- **Goroutine Profile**: Examines goroutine behavior
  - Shows currently running goroutines and their states
  - Helps identify goroutine leaks and deadlocks
  - Provides detailed stack traces for debugging

- **Block Profile**: Analyzes synchronization issues
  - Shows where goroutines block on synchronization primitives
  - Helps identify contention points and performance bottlenecks
  - Includes blocking duration statistics

- **Allocation Profile**: Examines memory allocation patterns
  - Shows memory allocation frequency by location
  - Helps identify memory churn and optimization opportunities
  - Includes both allocated and freed memory statistics

- **Thread Creation Profile**: Monitors OS thread creation
  - Shows where new OS threads are created
  - Helps identify thread leaks and excessive thread creation
  - Useful for analyzing thread pool behavior

## Profile View Modes

Each profile can be viewed in three different modes:

- **Flat View** (default): Shows direct values for each function
  - Displays the time/memory/etc. spent directly in each function
  - Excludes time spent in functions called by this function
  - Best for identifying specific hot spots in the code

- **Cumulative View**: Shows cumulative values including child functions
  - Includes time spent in the function and all functions it calls
  - Helps identify high-level bottlenecks in the call hierarchy
  - Useful for understanding the full impact of function calls

- **Graph View**: Shows the call graph relationship
  - Displays parent-child relationships between functions
  - Shows the top 5 children for each function
  - Helps understand the call flow and identify problematic paths

## Usage

### Basic Integration

```go
package main

import (
    "context"
    "log"
    pprofmcpagent "github.com/yudppp/pprof-mcp-agent"
)

func main() {
    // Create a context for controlling the server lifecycle
    ctx := context.Background()

    // Start the MCP server on port 1239
    if err := pprofmcpagent.ServeSSE(ctx, ":1239"); err != nil {
        log.Fatalf("Failed to start pprof MCP agent: %v", err)
    }
}
```

### Configuration

Each profile type supports the following configuration options:

- `limit`: Maximum number of locations to show in results (default: 100, min: 100, max: 10000)
- `view`: Profile view mode (`flat`, `cum`, or `graph`, default: `flat`)
- `duration`: Sampling duration for CPU profiles (default: 10 seconds)

## Features

- **Real-time Profiling**: Collect profiling data from running applications
- **Multiple View Modes**: Analyze data in flat, cumulative, or graph views
- **Aggregated Results**: View aggregated statistics for better analysis
- **Multiple Profile Types**: Comprehensive coverage of different performance aspects
- **Easy Integration**: Simple API for quick integration into existing applications
- **MCP Protocol**: Standard protocol for reliable data delivery
- **Configurable**: Adjustable parameters for different profiling needs

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your PR includes appropriate tests and documentation updates.

## License

[The MIT License (MIT)](https://github.com/yudppp/pprof-mcp-agent/blob/main/LICENSE)

## Author

yudppp

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
