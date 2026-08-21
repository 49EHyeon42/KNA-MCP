package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	klnimcp "github.com/49EHyeon42/KNA-MCP/internal/klni/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/service"
)

func addKlniTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return klnimcp.AddTools(server, klnimcp.UseCases{
		ScnmSearch: service.NewScnmSearchService(client),
	})
}
