package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSmplSearchUseCase defines the plant sample search use case.
type PlantSmplSearchUseCase interface {
	PlantSmplSearch(context.Context, application.PlantSmplSearchQuery) (application.PlantSmplSearchResult, error)
}
