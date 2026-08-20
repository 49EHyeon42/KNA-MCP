package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	kfnimcp "github.com/49EHyeon42/KNA-MCP/internal/kfni/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/service"
)

func addKfniTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return kfnimcp.AddTools(server, kfnimcp.UseCases{
		ScnmSearch: service.NewScnmSearchService(client),
		ScnmInfo:   service.NewScnmInfoService(client),
	})
}
