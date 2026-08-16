package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantNaturalizedListPort defines the outbound port for naturalized plant lists.
type PlantNaturalizedListPort interface {
	PlantNaturalizedList(context.Context, application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error)
}
