package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnIlstrSearchPort defines the outbound port for lichen pictorial book searches.
type AlchnIlstrSearchPort interface {
	AlchnIlstrSearch(context.Context, application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error)
}
