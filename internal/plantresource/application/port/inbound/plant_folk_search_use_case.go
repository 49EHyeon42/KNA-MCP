package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantFolkSearchUseCase defines the folk plant search use case.
type PlantFolkSearchUseCase interface {
	PlantFolkSearch(context.Context, application.PlantFolkSearchQuery) (application.PlantFolkSearchResult, error)
}
