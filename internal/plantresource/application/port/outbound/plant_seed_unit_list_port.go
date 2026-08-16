package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedUnitListPort defines the outbound port for plant seed unit lists.
type PlantSeedUnitListPort interface {
	PlantSeedUnitList(context.Context, application.PlantSeedUnitListQuery) (application.PlantSeedUnitListResult, error)
}
