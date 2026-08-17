package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	plantmstnsmcp "github.com/49EHyeon42/KNA-MCP/internal/plantmstns/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/service"
)

func addPlantMstnsTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return plantmstnsmcp.AddTools(server, plantmstnsmcp.UseCases{
		PlantMstnsList: service.NewPlantMstnsListService(client),
	})
}
