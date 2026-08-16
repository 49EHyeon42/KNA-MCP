package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	plantresourcemcp "github.com/49EHyeon42/KNA-MCP/internal/plantresource/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/service"
)

func addPlantResourceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return plantresourcemcp.AddTools(server, plantresourcemcp.UseCases{
		PlantPilbkSearch:     service.NewPlantPilbkSearchService(client),
		PlantPilbkInfo:       service.NewPlantPilbkInfoService(client),
		PlantSmplSearch:      service.NewPlantSmplSearchService(client),
		PlantSmplUnitList:    service.NewPlantSmplUnitListService(client),
		PlantSeedSearch:      service.NewPlantSeedSearchService(client),
		PlantSeedUnitList:    service.NewPlantSeedUnitListService(client),
		PlantSeedGrmntList:   service.NewPlantSeedGrmntListService(client),
		PlantFolkSearch:      service.NewPlantFolkSearchService(client),
		PlantFolkAreaList:    service.NewPlantFolkAreaListService(client),
		PlantNaturalizedList: service.NewPlantNaturalizedListService(client),
		PlantRareList:        service.NewPlantRareListService(client),
		PlantSpcltList:       service.NewPlantSpcltListService(client),
		PlantWordList:        service.NewPlantWordListService(client),
	})
}
