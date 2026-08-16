package main

import (
	"context"
	"log"
	"os"

	"github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func main() {
	server := mcpstdio.NewServer()
	serviceKey := os.Getenv("DATA_GO_KR_SERVICE_KEY")
	if err := addPlantResourceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addPlantMstnsTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}

	if err := mcpstdio.Run(context.Background(), server); err != nil {
		log.Fatal(err)
	}
}
