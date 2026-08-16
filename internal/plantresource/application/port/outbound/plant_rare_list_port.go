package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantRareListPort defines the outbound port for rare plant lists.
type PlantRareListPort interface {
	PlantRareList(context.Context, application.PlantRareListQuery) (application.PlantRareListResult, error)
}
