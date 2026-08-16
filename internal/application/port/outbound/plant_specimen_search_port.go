package outbound

import (
	"context"

	"kna-mcp/internal/application"
)

// PlantSpecimenSearchPort defines the outbound port for plant specimen searches.
type PlantSpecimenSearchPort interface {
	PlantSpecimenSearch(context.Context, application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error)
}
