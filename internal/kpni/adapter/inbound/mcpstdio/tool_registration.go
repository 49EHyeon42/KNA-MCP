package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
)

// UseCases contains the national standard plant list use cases exposed as MCP tools.
type UseCases struct {
	ScnmSearch inbound.ScnmSearchUseCase
}

// AddTools adds all national standard plant list tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.ScnmSearch == nil {
		return errors.New("scnmSearch use case is required")
	}

	addScnmSearchTool(server, useCases.ScnmSearch)
	return nil
}
