package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/port/inbound"
)

// UseCases contains the biological information center use cases exposed as MCP tools.
type UseCases struct {
	RelatedSiteList inbound.RelatedSiteListUseCase
}

// AddTools adds all biological information center tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.RelatedSiteList == nil {
		return errors.New("relatedSiteList use case is required")
	}

	addRelatedSiteListTool(server, useCases.RelatedSiteList)
	return nil
}
