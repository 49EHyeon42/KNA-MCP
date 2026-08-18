// Package mcpstdio runs the MCP server over standard input and output.
package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewServer creates the KNA MCP server.
func NewServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{Name: "kna-mcp", Version: "0.5.0"}, nil)
}

// Run serves MCP over standard input and output.
func Run(ctx context.Context, server *mcp.Server) error {
	return server.Run(ctx, &mcp.StdioTransport{})
}
