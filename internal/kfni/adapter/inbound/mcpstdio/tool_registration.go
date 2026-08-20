package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/port/inbound"
)

// UseCases contains the national standard fungi list use cases exposed as MCP tools.
type UseCases struct {
	ScnmSearch inbound.ScnmSearchUseCase
	ScnmInfo   inbound.ScnmInfoUseCase
}

// AddTools adds all national standard fungi list tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.ScnmSearch == nil {
		return errors.New("scnmSearch use case is required")
	}
	if useCases.ScnmInfo == nil {
		return errors.New("scnmInfo use case is required")
	}

	addScnmSearchTool(server, useCases.ScnmSearch)
	addScnmInfoTool(server, useCases.ScnmInfo)
	return nil
}
