package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsPilbkSearchPort defines the outbound port for fungi pictorial book searches.
type FngsPilbkSearchPort interface {
	FngsPilbkSearch(context.Context, application.FngsPilbkSearchQuery) (application.FngsPilbkSearchResult, error)
}
