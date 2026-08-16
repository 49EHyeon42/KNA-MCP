package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/application/port/inbound"
)

// UseCases contains the application use cases exposed as MCP tools.
type UseCases struct {
	PlantPictorialBookSearch      inbound.PlantPictorialBookSearchUseCase
	PlantPictorialBookInformation inbound.PlantPictorialBookInformationUseCase
	PlantSpecimenSearch           inbound.PlantSpecimenSearchUseCase
}

// Run serves MCP over standard input and output.
func Run(ctx context.Context, useCases UseCases) error {
	return newServer(useCases).Run(ctx, &mcp.StdioTransport{})
}

func newServer(useCases UseCases) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "kna-mcp", Version: "0.1.0"}, nil)
	if useCases.PlantPictorialBookSearch != nil {
		addPlantPictorialBookSearchTool(server, useCases.PlantPictorialBookSearch)
	}
	if useCases.PlantPictorialBookInformation != nil {
		addPlantPictorialBookInformationTool(server, useCases.PlantPictorialBookInformation)
	}
	if useCases.PlantSpecimenSearch != nil {
		addPlantSpecimenSearchTool(server, useCases.PlantSpecimenSearch)
	}
	return server
}
