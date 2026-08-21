package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
)

// ScnmSearchPort defines the outbound port for lichen scientific name lists.
type ScnmSearchPort interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
