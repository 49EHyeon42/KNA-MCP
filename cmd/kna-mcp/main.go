package main

import (
	"context"
	"log"
	"os"

	"github.com/49EHyeon42/KNA-MCP/internal/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/application/service"
)

func main() {
	client, err := kna.NewClient(os.Getenv("DATA_GO_KR_SERVICE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	if err := mcpstdio.Run(context.Background(), mcpstdio.UseCases{
		PlantPictorialBookSearch:      service.NewPlantPictorialBookSearchService(client),
		PlantPictorialBookInformation: service.NewPlantPictorialBookInformationService(client),
		PlantSampleSearch:             service.NewPlantSampleSearchService(client),
	}); err != nil {
		log.Fatal(err)
	}
}
