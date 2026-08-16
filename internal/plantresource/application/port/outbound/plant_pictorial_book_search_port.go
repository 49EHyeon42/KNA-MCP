package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPictorialBookSearchPort defines the outbound port for plant pictorial book searches.
type PlantPictorialBookSearchPort interface {
	PlantPictorialBookSearch(context.Context, application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error)
}
