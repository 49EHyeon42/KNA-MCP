package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/inbound"
)

// UseCases contains the plant miniature use cases exposed as MCP tools.
type UseCases struct {
	PlantMstnsList inbound.PlantMstnsListUseCase
}

// AddTools adds all plant miniature tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.PlantMstnsList == nil {
		return errors.New("plantMstnsList use case is required")
	}

	addPlantMstnsListTool(server, useCases.PlantMstnsList)
	return nil
}
