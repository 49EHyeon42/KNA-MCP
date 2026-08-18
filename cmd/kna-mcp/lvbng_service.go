package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	lvbngservicemcp "github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/service"
)

func addLvbngServiceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return lvbngservicemcp.AddTools(server, lvbngservicemcp.UseCases{
		RelatedSiteList: service.NewRelatedSiteListService(client),
	})
}
