package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSmplUnitListPort defines the outbound port for plant specimen detail lists.
type PlantSmplUnitListPort interface {
	PlantSmplUnitList(context.Context, application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error)
}
