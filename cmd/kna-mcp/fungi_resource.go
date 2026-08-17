package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	fungiresourcemcp "github.com/49EHyeon42/KNA-MCP/internal/fungiresource/adapter/inbound/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/adapter/outbound/kna"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/service"
)

func addFungiResourceTools(server *mcp.Server, serviceKey string) error {
	client, err := kna.NewClient(serviceKey)
	if err != nil {
		return err
	}

	return fungiresourcemcp.AddTools(server, fungiresourcemcp.UseCases{
		FngsPilbkSearch: service.NewFngsPilbkSearchService(client),
		FngsPilbkInfo:   service.NewFngsPilbkInfoService(client),
		FngsSmplSearch:  service.NewFngsSmplSearchService(client),
	})
}
