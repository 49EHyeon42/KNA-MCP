package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

// PlantSpecimenSearchPort defines the outbound port for plant specimen searches.
type PlantSpecimenSearchPort interface {
	PlantSpecimenSearch(context.Context, application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error)
}
