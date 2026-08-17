package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	kpnimcp "github.com/49EHyeon42/KNA-MCP/internal/kpni/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/service"
)

func addKpniTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return kpnimcp.AddTools(server, kpnimcp.UseCases{
		ScnmSearch: service.NewScnmSearchService(client),
		ScnmInfo:   service.NewScnmInfoService(client),
	})
}
