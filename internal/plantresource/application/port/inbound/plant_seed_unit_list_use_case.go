package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedUnitListUseCase defines the plant seed unit list use case.
type PlantSeedUnitListUseCase interface {
	PlantSeedUnitList(context.Context, application.PlantSeedUnitListQuery) (application.PlantSeedUnitListResult, error)
}
