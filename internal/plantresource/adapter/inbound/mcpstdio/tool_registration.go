package mcpstdio

import (
	"errors"

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

// AddTools adds all plant resource tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	switch {
	case useCases.PlantPilbkSearch == nil:
		return errors.New("plantPilbkSearch use case is required")
	case useCases.PlantPilbkInfo == nil:
		return errors.New("plantPilbkInfo use case is required")
	case useCases.PlantSmplSearch == nil:
		return errors.New("plantSmplSearch use case is required")
	case useCases.PlantSmplUnitList == nil:
		return errors.New("plantSmplUnitList use case is required")
	case useCases.PlantSeedSearch == nil:
		return errors.New("plantSeedSearch use case is required")
	case useCases.PlantSeedUnitList == nil:
		return errors.New("plantSeedUnitList use case is required")
	case useCases.PlantSeedGrmntList == nil:
		return errors.New("plantSeedGrmntList use case is required")
	case useCases.PlantFolkSearch == nil:
		return errors.New("plantFolkSearch use case is required")
	case useCases.PlantFolkAreaList == nil:
		return errors.New("plantFolkAreaList use case is required")
	case useCases.PlantNaturalizedList == nil:
		return errors.New("plantNaturalizedList use case is required")
	case useCases.PlantRareList == nil:
		return errors.New("plantRareList use case is required")
	case useCases.PlantSpcltList == nil:
		return errors.New("plantSpcltList use case is required")
	case useCases.PlantWordList == nil:
		return errors.New("plantWordList use case is required")
	}

	addPlantPilbkSearchTool(server, useCases.PlantPilbkSearch)
	addPlantPilbkInfoTool(server, useCases.PlantPilbkInfo)
	addPlantSmplSearchTool(server, useCases.PlantSmplSearch)
	addPlantSmplUnitListTool(server, useCases.PlantSmplUnitList)
	addPlantSeedSearchTool(server, useCases.PlantSeedSearch)
	addPlantSeedUnitListTool(server, useCases.PlantSeedUnitList)
	addPlantSeedGrmntListTool(server, useCases.PlantSeedGrmntList)
	addPlantFolkSearchTool(server, useCases.PlantFolkSearch)
	addPlantFolkAreaListTool(server, useCases.PlantFolkAreaList)
	addPlantNaturalizedListTool(server, useCases.PlantNaturalizedList)
	addPlantRareListTool(server, useCases.PlantRareList)
	addPlantSpcltListTool(server, useCases.PlantSpcltList)
	addPlantWordListTool(server, useCases.PlantWordList)

	return nil
}
