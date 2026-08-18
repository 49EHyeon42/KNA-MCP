package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
)

// UseCases contains the lichen service use cases exposed as MCP tools.
type UseCases struct {
	AlchnIlstrSearch inbound.AlchnIlstrSearchUseCase
	AlchnIlstrInfo   inbound.AlchnIlstrInfoUseCase
	AlchnSpcmSearch  inbound.AlchnSpcmSearchUseCase
	AlchnSpcmInfo    inbound.AlchnSpcmInfoUseCase
}

// AddTools adds all lichen service tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.AlchnIlstrSearch == nil {
		return errors.New("alchnIlstrSearch use case is required")
	}
	if useCases.AlchnIlstrInfo == nil {
		return errors.New("alchnIlstrInfo use case is required")
	}
	if useCases.AlchnSpcmSearch == nil {
		return errors.New("alchnSpcmSearch use case is required")
	}
	if useCases.AlchnSpcmInfo == nil {
		return errors.New("alchnSpcmInfo use case is required")
	}

	addAlchnIlstrSearchTool(server, useCases.AlchnIlstrSearch)
	addAlchnIlstrInfoTool(server, useCases.AlchnIlstrInfo)
	addAlchnSpcmSearchTool(server, useCases.AlchnSpcmSearch)
	addAlchnSpcmInfoTool(server, useCases.AlchnSpcmInfo)
	return nil
}
