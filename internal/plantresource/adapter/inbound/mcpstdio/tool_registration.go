package mcpstdio

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

// UseCases contains the plant resource use cases exposed as MCP tools.
type UseCases struct {
	PlantPilbkSearch inbound.PlantPilbkSearchUseCase
	PlantPilbkInfo   inbound.PlantPilbkInfoUseCase
	PlantSmplSearch  inbound.PlantSmplSearchUseCase
}

// AddTools adds the plant resource tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) {
	if useCases.PlantPilbkSearch != nil {
		addPlantPilbkSearchTool(server, useCases.PlantPilbkSearch)
	}
	if useCases.PlantPilbkInfo != nil {
		addPlantPilbkInfoTool(server, useCases.PlantPilbkInfo)
	}
	if useCases.PlantSmplSearch != nil {
		addPlantSmplSearchTool(server, useCases.PlantSmplSearch)
	}
}
