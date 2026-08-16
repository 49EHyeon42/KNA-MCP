package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantRareListUseCase defines the rare plant list use case.
type PlantRareListUseCase interface {
	PlantRareList(context.Context, application.PlantRareListQuery) (application.PlantRareListResult, error)
}
