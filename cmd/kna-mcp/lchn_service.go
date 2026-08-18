package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	lchnservicemcp "github.com/49EHyeon42/KNA-MCP/internal/lchnservice/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/service"
)

func addLchnServiceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return lchnservicemcp.AddTools(server, lchnservicemcp.UseCases{
		AlchnIlstrSearch: service.NewAlchnIlstrSearchService(client),
		AlchnIlstrInfo:   service.NewAlchnIlstrInfoService(client),
	})
}
