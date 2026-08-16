package mcpstdio

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

// UseCases contains the plant resource use cases exposed as MCP tools.
type UseCases struct {
	PlantPictorialBookSearch      inbound.PlantPictorialBookSearchUseCase
	PlantPictorialBookInformation inbound.PlantPictorialBookInformationUseCase
	PlantSampleSearch             inbound.PlantSampleSearchUseCase
}

// AddTools adds the plant resource tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) {
	if useCases.PlantPictorialBookSearch != nil {
		addPlantPictorialBookSearchTool(server, useCases.PlantPictorialBookSearch)
	}
	if useCases.PlantPictorialBookInformation != nil {
		addPlantPictorialBookInformationTool(server, useCases.PlantPictorialBookInformation)
	}
	if useCases.PlantSampleSearch != nil {
		addPlantSampleSearchTool(server, useCases.PlantSampleSearch)
	}
}
