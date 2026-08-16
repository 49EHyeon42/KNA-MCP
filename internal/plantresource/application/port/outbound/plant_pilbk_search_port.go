package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPilbkSearchPort defines the outbound port for plant pictorial book searches.
type PlantPilbkSearchPort interface {
	PlantPilbkSearch(context.Context, application.PlantPilbkSearchQuery) (application.PlantPilbkSearchResult, error)
}
