package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantFolkSearchPort defines the outbound port for folk plant searches.
type PlantFolkSearchPort interface {
	PlantFolkSearch(context.Context, application.PlantFolkSearchQuery) (application.PlantFolkSearchResult, error)
}
