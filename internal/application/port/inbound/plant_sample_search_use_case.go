package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

// PlantSampleSearchUseCase defines the plant sample search use case.
type PlantSampleSearchUseCase interface {
	PlantSampleSearch(context.Context, application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error)
}
