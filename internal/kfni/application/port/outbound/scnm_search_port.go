package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
)

// ScnmSearchPort defines the outbound port for fungi scientific name lists.
type ScnmSearchPort interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
