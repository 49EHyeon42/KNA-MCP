package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantNaturalizedListUseCase defines the naturalized plant list use case.
type PlantNaturalizedListUseCase interface {
	PlantNaturalizedList(context.Context, application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error)
}
