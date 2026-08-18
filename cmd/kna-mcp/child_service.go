package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	childservicemcp "github.com/49EHyeon42/KNA-MCP/internal/childservice/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/service"
)

func addChildServiceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return childservicemcp.AddTools(server, childservicemcp.UseCases{
		ChildPilbkSearch: service.NewChildPilbkSearchService(client),
		ChildPilbkInfo:   service.NewChildPilbkInfoService(client),
	})
}
