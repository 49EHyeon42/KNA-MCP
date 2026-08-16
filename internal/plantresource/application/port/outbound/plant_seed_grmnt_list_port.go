package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSeedGrmntListPort defines the outbound port for plant seed germination lists.
type PlantSeedGrmntListPort interface {
	PlantSeedGrmntList(context.Context, application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error)
}
