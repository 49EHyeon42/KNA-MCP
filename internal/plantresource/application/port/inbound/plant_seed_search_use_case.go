package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedSearchUseCase defines the plant seed search use case.
type PlantSeedSearchUseCase interface {
	PlantSeedSearch(context.Context, application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error)
}
