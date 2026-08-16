package inbound

import (
	"context"

	"kna-mcp/internal/application"
)

// PlantSpecimenSearchUseCase defines the plant specimen search use case.
type PlantSpecimenSearchUseCase interface {
	PlantSpecimenSearch(context.Context, application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error)
}
