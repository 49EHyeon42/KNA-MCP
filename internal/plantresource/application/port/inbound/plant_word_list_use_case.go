package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantWordListUseCase defines the plant word list use case.
type PlantWordListUseCase interface {
	PlantWordList(context.Context, application.PlantWordListQuery) (application.PlantWordListResult, error)
}
