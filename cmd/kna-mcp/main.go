package main

import (
	"context"
	"log"
	"os"

	"kna-mcp/internal/adapter/inbound/mcpstdio"
	"kna-mcp/internal/adapter/outbound/kna"
	"kna-mcp/internal/application/service"
)

func main() {
	client, err := kna.NewClient(os.Getenv("DATA_GO_KR_SERVICE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	if err := mcpstdio.Run(context.Background(), mcpstdio.UseCases{
		PlantPictorialBookSearch: service.NewPlantPictorialBookSearchService(client),
		PlantSpecimenSearch:      service.NewPlantSpecimenSearchService(client),
	}); err != nil {
		log.Fatal(err)
	}
}
