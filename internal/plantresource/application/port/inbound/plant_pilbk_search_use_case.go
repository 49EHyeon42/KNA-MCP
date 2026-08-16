package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPilbkSearchUseCase defines the plant pictorial book search use case.
type PlantPilbkSearchUseCase interface {
	PlantPilbkSearch(context.Context, application.PlantPilbkSearchQuery) (application.PlantPilbkSearchResult, error)
}
