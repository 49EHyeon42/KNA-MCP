package main

import (
	"context"
	"log"
	"os"

	"github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func main() {
	server := mcpstdio.NewServer()
	if err := addPlantResourceTools(server, os.Getenv("DATA_GO_KR_SERVICE_KEY")); err != nil {
		log.Fatal(err)
	}

	if err := mcpstdio.Run(context.Background(), server); err != nil {
		log.Fatal(err)
	}
}
