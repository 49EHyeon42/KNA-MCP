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
	if err := addKpniTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addKiniTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addKfniTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addKlniTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addInsectResourceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addFungiResourceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addChildServiceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addLvbngServiceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addLchnServiceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}
	if err := addOldPlantServiceTools(server, serviceKey); err != nil {
		log.Fatal(err)
	}

	if err := mcpstdio.Run(context.Background(), server); err != nil {
		log.Fatal(err)
	}
}
