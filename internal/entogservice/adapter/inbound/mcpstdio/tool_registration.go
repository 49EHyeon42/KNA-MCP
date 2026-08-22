package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

// UseCases contains the entognath service use cases exposed as MCP tools.
type UseCases struct {
	EntogIlstrSearch inbound.EntogIlstrSearchUseCase
}

// AddTools adds all entognath service tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.EntogIlstrSearch == nil {
		return errors.New("entogIlstrSearch use case is required")
	}

	addEntogIlstrSearchTool(server, useCases.EntogIlstrSearch)
	return nil
}
