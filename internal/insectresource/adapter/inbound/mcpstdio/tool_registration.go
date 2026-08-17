package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

// UseCases contains the insect resource use cases exposed as MCP tools.
type UseCases struct {
	InsectPilbkSearch inbound.InsectPilbkSearchUseCase
	InsectPilbkInfo   inbound.InsectPilbkInfoUseCase
}

// AddTools adds all insect resource tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	switch {
	case useCases.InsectPilbkSearch == nil:
		return errors.New("insectPilbkSearch use case is required")
	case useCases.InsectPilbkInfo == nil:
		return errors.New("insectPilbkInfo use case is required")
	}

	addInsectPilbkSearchTool(server, useCases.InsectPilbkSearch)
	addInsectPilbkInfoTool(server, useCases.InsectPilbkInfo)
	return nil
}
