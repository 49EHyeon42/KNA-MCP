package mcpstdio

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

// UseCases contains the plant resource use cases exposed as MCP tools.
type UseCases struct {
	PlantPilbkSearch     inbound.PlantPilbkSearchUseCase
	PlantPilbkInfo       inbound.PlantPilbkInfoUseCase
	PlantSmplSearch      inbound.PlantSmplSearchUseCase
	PlantSmplUnitList    inbound.PlantSmplUnitListUseCase
	PlantSeedSearch      inbound.PlantSeedSearchUseCase
	PlantSeedUnitList    inbound.PlantSeedUnitListUseCase
	PlantSeedGrmntList   inbound.PlantSeedGrmntListUseCase
	PlantFolkSearch      inbound.PlantFolkSearchUseCase
	PlantFolkAreaList    inbound.PlantFolkAreaListUseCase
	PlantNaturalizedList inbound.PlantNaturalizedListUseCase
	PlantRareList        inbound.PlantRareListUseCase
	PlantSpcltList       inbound.PlantSpcltListUseCase
	PlantWordList        inbound.PlantWordListUseCase
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
	if useCases.PlantSmplUnitList != nil {
		addPlantSmplUnitListTool(server, useCases.PlantSmplUnitList)
	}
	if useCases.PlantSeedSearch != nil {
		addPlantSeedSearchTool(server, useCases.PlantSeedSearch)
	}
	if useCases.PlantSeedUnitList != nil {
		addPlantSeedUnitListTool(server, useCases.PlantSeedUnitList)
	}
	if useCases.PlantSeedGrmntList != nil {
		addPlantSeedGrmntListTool(server, useCases.PlantSeedGrmntList)
	}
	if useCases.PlantFolkSearch != nil {
		addPlantFolkSearchTool(server, useCases.PlantFolkSearch)
	}
	if useCases.PlantFolkAreaList != nil {
		addPlantFolkAreaListTool(server, useCases.PlantFolkAreaList)
	}
	if useCases.PlantNaturalizedList != nil {
		addPlantNaturalizedListTool(server, useCases.PlantNaturalizedList)
	}
	if useCases.PlantRareList != nil {
		addPlantRareListTool(server, useCases.PlantRareList)
	}
	if useCases.PlantSpcltList != nil {
		addPlantSpcltListTool(server, useCases.PlantSpcltList)
	}
	if useCases.PlantWordList != nil {
		addPlantWordListTool(server, useCases.PlantWordList)
	}
}
