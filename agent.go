package pprofmcpagent

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime/pprof"
	"time"

	"github.com/felixge/pprofutils/v2/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ServeMCPServer(ctx context.Context, port string) error {
	s := server.NewMCPServer(
		"pprof server",
		"1.0.0",
	)

	// Add tool
	tool := NewTool()

	// Add tool handler
	s.AddTool(tool, Handler)

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

	// Add debug logging
	log.Printf("Starting SSE server on port %s\n", port)

	if err := sses.Start(port); err != nil {
		log.Printf("SSE server error: %v\n", err)
		return fmt.Errorf("SSE server error: %w", err)
	}

	return nil
}

func NewTool() mcp.Tool {
	return mcp.NewTool("pprof server",
		mcp.WithDescription("Output pprof data"),
		mcp.WithString(
			"profile",
			mcp.Description(
				"The type of profile to output (e.g., heap, goroutine, threadcreate, block, allocs, cpu)",
			),
			mcp.Required(),
			mcp.Enum("heap", "goroutine", "threadcreate", "block", "allocs", "cpu"),
		),
		mcp.WithNumber(
			"duration",
			mcp.Description("Duration in seconds(number) to collect the profile (only for CPU profile)"),
			mcp.DefaultNumber(10),
			mcp.Min(1),
			mcp.Max(60),
		),
	)
}

func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profileType, ok := request.Params.Arguments["profile"].(string)
	if !ok {
		return handleMCPError(fmt.Errorf("profile argument is required and must be a string")), nil
	}
	if profileType == "cpu" {
		duration := 10
		if d, ok := request.Params.Arguments["duration"].(float64); ok {
			duration = int(d)
		}

		var buf bytes.Buffer
		if err := pprof.StartCPUProfile(&buf); err != nil {
			return handleMCPError(err), nil
		}
		time.Sleep(time.Duration(duration) * time.Second)
		pprof.StopCPUProfile()

		data := bytes.NewBuffer(nil)
		pprofutils := &utils.JSON{Input: buf.Bytes(), Output: data}
		if err := pprofutils.Execute(ctx); err != nil {
			return handleMCPError(err), nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(data.String()),
			},
		}, nil
	}

	profile := pprof.Lookup(profileType)
	if profile == nil {
		return handleMCPError(fmt.Errorf("profile %v not found", profileType)), nil

	}
	data, err := parsePofileData(ctx, profile)
	if err != nil {
		return handleMCPError(err), nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(data),
		},
	}, nil
}

func parsePofileData(ctx context.Context, profile *pprof.Profile) (string, error) {
	input := bytes.NewBuffer(nil)
	output := bytes.NewBuffer(nil)
	if err := profile.WriteTo(input, 0); err != nil {
		return "", err
	}
	pprofutils := &utils.JSON{Input: input.Bytes(), Output: output}
	if err := pprofutils.Execute(ctx); err != nil {
		return "", err
	}
	return output.String(), nil
}

func handleMCPError(err error) *mcp.CallToolResult {
	log.Println(err)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(err.Error()),
		},
		IsError: true,
	}
}
