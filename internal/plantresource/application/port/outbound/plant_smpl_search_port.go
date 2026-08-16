package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantSmplSearchPort defines the outbound port for plant sample searches.
type PlantSmplSearchPort interface {
	PlantSmplSearch(context.Context, application.PlantSmplSearchQuery) (application.PlantSmplSearchResult, error)
}
