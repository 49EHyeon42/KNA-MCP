package mcpstdio

import (
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

// UseCases contains the entognath service use cases exposed as MCP tools.
type UseCases struct {
	EntogIlstrSearch inbound.EntogIlstrSearchUseCase
	EntogIlstrInfo   inbound.EntogIlstrInfoUseCase
	EntogSpcmSearch  inbound.EntogSpcmSearchUseCase
}

// AddTools adds all entognath service tools to an MCP server.
func AddTools(server *mcp.Server, useCases UseCases) error {
	if useCases.EntogIlstrSearch == nil {
		return errors.New("entogIlstrSearch use case is required")
	}
	if useCases.EntogIlstrInfo == nil {
		return errors.New("entogIlstrInfo use case is required")
	}
	if useCases.EntogSpcmSearch == nil {
		return errors.New("entogSpcmSearch use case is required")
	}

	addEntogIlstrSearchTool(server, useCases.EntogIlstrSearch)
	addEntogIlstrInfoTool(server, useCases.EntogIlstrInfo)
	addEntogSpcmSearchTool(server, useCases.EntogSpcmSearch)
	return nil
}
