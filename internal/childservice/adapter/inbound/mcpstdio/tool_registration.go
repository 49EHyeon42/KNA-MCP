package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/inbound"
)

// UseCases contains the child pictorial book use cases exposed as MCP tools.
type UseCases struct {
	ChildPilbkSearch inbound.ChildPilbkSearchUseCase
	ChildPilbkInfo   inbound.ChildPilbkInfoUseCase
}

// AddTools adds all child pictorial book tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.ChildPilbkSearch == nil {
		return errors.New("childPilbkSearch use case is required")
	}
	if useCases.ChildPilbkInfo == nil {
		return errors.New("childPilbkInfo use case is required")
	}

	addChildPilbkSearchTool(server, useCases.ChildPilbkSearch)
	addChildPilbkInfoTool(server, useCases.ChildPilbkInfo)
	return nil
}
