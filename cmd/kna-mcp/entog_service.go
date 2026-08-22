package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	entogservicemcp "github.com/49EHyeon42/KNA-MCP/internal/entogservice/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/service"
)

func addEntogServiceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return entogservicemcp.AddTools(server, entogservicemcp.UseCases{
		EntogIlstrSearch: service.NewEntogIlstrSearchService(client),
		EntogIlstrInfo:   service.NewEntogIlstrInfoService(client),
		EntogSpcmSearch:  service.NewEntogSpcmSearchService(client),
	})
}
