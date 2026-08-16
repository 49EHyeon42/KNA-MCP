package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSmplUnitListUseCase defines the plant specimen detail list use case.
type PlantSmplUnitListUseCase interface {
	PlantSmplUnitList(context.Context, application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error)
}
