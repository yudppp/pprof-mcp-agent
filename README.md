# pprof-mcp-agent

A Go agent that provides pprof profiling data through an MCP (Mission Control Protocol) server.

## Installation

```bash
go get github.com/yudppp/pprof-mcp-agent
```

## Overview

pprof-mcp-agent is a tool that makes it easy to collect and analyze performance profiling data from Go applications. Using the MCP protocol, it provides various profile data types including:

- CPU profiles
- Heap profiles
- Goroutine profiles
- Thread creation profiles
- Block profiles
- Allocation profiles

## Features

- Profile data delivery via MCP protocol
- Easy integration and configuration
- Real-time profiling data collection
- Customizable profiling duration (for CPU profiles)
- Graceful shutdown support

## Usage

```go
package main

import (
    "context"
    pprofmcpagent "github.com/yudppp/pprof-mcp-agent"
)

func main() {
    ctx := context.Background()
    err := pprofmcpagent.ServeMCPServer(ctx, ":1239")
    if err != nil {
        // Handle error
    }
}
```

## Profile Types

The following profile types are supported:

- `cpu`: CPU profile (default 10-second collection)
- `heap`: Heap memory profile
- `goroutine`: Goroutine profile
- `threadcreate`: Thread creation profile
- `block`: Block profile
- `allocs`: Allocation profile

## Configuration Options

### CPU Profile

- `duration`: Profile collection time (seconds)
  - Minimum: 1 second
  - Maximum: 60 seconds
  - Default: 10 seconds

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

## Author

yudppp
