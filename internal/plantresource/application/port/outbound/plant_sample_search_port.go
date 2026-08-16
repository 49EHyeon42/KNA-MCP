package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSampleSearchPort defines the outbound port for plant sample searches.
type PlantSampleSearchPort interface {
	PlantSampleSearch(context.Context, application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error)
}
