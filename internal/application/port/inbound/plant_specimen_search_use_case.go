package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

// PlantSpecimenSearchUseCase defines the plant specimen search use case.
type PlantSpecimenSearchUseCase interface {
	PlantSpecimenSearch(context.Context, application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error)
}
