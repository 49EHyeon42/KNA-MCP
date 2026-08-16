package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSpcltListUseCase defines the endemic plant list use case.
type PlantSpcltListUseCase interface {
	PlantSpcltList(context.Context, application.PlantSpcltListQuery) (application.PlantSpcltListResult, error)
}
