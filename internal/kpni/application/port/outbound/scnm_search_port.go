package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

// ScnmSearchPort defines the outbound port for scientific name lists.
type ScnmSearchPort interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
