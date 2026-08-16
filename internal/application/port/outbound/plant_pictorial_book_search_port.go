package outbound

import (
	"context"

	"kna-mcp/internal/application"
)

// PlantPictorialBookSearchPort defines the outbound port for plant pictorial book searches.
type PlantPictorialBookSearchPort interface {
	PlantPictorialBookSearch(context.Context, application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error)
}
