package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	kinimcp "github.com/49EHyeon42/KNA-MCP/internal/kini/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/application/service"
)

func addKiniTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return kinimcp.AddTools(server, kinimcp.UseCases{
		ScnmSearch: service.NewScnmSearchService(client),
		ScnmInfo:   service.NewScnmInfoService(client),
	})
}
