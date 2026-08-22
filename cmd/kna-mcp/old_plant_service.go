package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	oldplantservicemcp "github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/service"
)

func addOldPlantServiceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return oldplantservicemcp.AddTools(server, oldplantservicemcp.UseCases{
		OldSpcmSearch: service.NewOldSpcmSearchService(client),
	})
}
