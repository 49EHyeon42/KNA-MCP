package main

import (
	"context"
	"log"
	"os"

	"github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
	plantresourcemcp "github.com/49EHyeon42/KNA-MCP/internal/plantresource/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/service"
)

func main() {
	client, err := kna.NewClient(os.Getenv("DATA_GO_KR_SERVICE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	server := mcpstdio.NewServer()
	plantresourcemcp.AddTools(server, plantresourcemcp.UseCases{
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
	})

	if err := mcpstdio.Run(context.Background(), server); err != nil {
		log.Fatal(err)
	}
}
