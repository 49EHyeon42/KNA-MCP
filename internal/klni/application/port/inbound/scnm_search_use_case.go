package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
)

// ScnmSearchUseCase defines the lichen scientific name list use case.
type ScnmSearchUseCase interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
