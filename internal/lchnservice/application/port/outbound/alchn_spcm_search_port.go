package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnSpcmSearchPort defines the outbound port for lichen specimen searches.
type AlchnSpcmSearchPort interface {
	AlchnSpcmSearch(context.Context, application.AlchnSpcmSearchQuery) (application.AlchnSpcmSearchResult, error)
}
