package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"kna-mcp/internal/application/port/inbound"
)

// Run serves MCP over standard input and output.
func Run(ctx context.Context, useCase inbound.PlantPictorialBookSearchUseCase) error {
	return newServer(useCase).Run(ctx, &mcp.StdioTransport{})
}

func newServer(useCase inbound.PlantPictorialBookSearchUseCase) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "kna-mcp", Version: "0.1.0"}, nil)
	addPlantPictorialBookSearchTool(server, useCase)
	return server
}
