package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/port/inbound"
)

// UseCases contains the old plant service use cases exposed as MCP tools.
type UseCases struct {
	OldSpcmSearch inbound.OldSpcmSearchUseCase
}

// AddTools adds all old plant service tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.OldSpcmSearch == nil {
		return errors.New("oldSpcmSearch use case is required")
	}

	addOldSpcmSearchTool(server, useCases.OldSpcmSearch)
	return nil
}
