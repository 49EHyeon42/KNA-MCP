package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedGrmntListUseCase defines the plant seed germination list use case.
type PlantSeedGrmntListUseCase interface {
	PlantSeedGrmntList(context.Context, application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error)
}
