package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	insectresourcemcp "github.com/49EHyeon42/KNA-MCP/internal/insectresource/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/service"
)

func addInsectResourceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return insectresourcemcp.AddTools(server, insectresourcemcp.UseCases{
		InsectPilbkSearch:  service.NewInsectPilbkSearchService(client),
		InsectPilbkInfo:    service.NewInsectPilbkInfoService(client),
		InsectPrtctList:    service.NewInsectPrtctListService(client),
		InsectSmplSearch:   service.NewInsectSmplSearchService(client),
		InsectSmplUnitList: service.NewInsectSmplUnitListService(client),
	})
}
