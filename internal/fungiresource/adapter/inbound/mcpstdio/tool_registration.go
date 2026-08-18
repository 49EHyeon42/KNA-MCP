package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
)

// UseCases contains the fungi resource use cases exposed as MCP tools.
type UseCases struct {
	FngsPilbkSearch  inbound.FngsPilbkSearchUseCase
	FngsPilbkInfo    inbound.FngsPilbkInfoUseCase
	FngsSmplSearch   inbound.FngsSmplSearchUseCase
	FngsSmplUnitList inbound.FngsSmplUnitListUseCase
}

// AddTools adds all fungi resource tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.FngsPilbkSearch == nil {
		return errors.New("fngsPilbkSearch use case is required")
	}
	if useCases.FngsPilbkInfo == nil {
		return errors.New("fngsPilbkInfo use case is required")
	}
	if useCases.FngsSmplSearch == nil {
		return errors.New("fngsSmplSearch use case is required")
	}
	if useCases.FngsSmplUnitList == nil {
		return errors.New("fngsSmplUnitList use case is required")
	}

	addFngsPilbkSearchTool(server, useCases.FngsPilbkSearch)
	addFngsPilbkInfoTool(server, useCases.FngsPilbkInfo)
	addFngsSmplSearchTool(server, useCases.FngsSmplSearch)
	addFngsSmplUnitListTool(server, useCases.FngsSmplUnitList)
	return nil
}
