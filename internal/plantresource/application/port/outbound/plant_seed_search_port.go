package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedSearchPort defines the outbound port for plant seed searches.
type PlantSeedSearchPort interface {
	PlantSeedSearch(context.Context, application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error)
}
