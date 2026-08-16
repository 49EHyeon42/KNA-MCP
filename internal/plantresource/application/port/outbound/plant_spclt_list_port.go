package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSpcltListPort defines the outbound port for endemic plant lists.
type PlantSpcltListPort interface {
	PlantSpcltList(context.Context, application.PlantSpcltListQuery) (application.PlantSpcltListResult, error)
}
