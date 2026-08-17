package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

// ScnmSearchUseCase defines the scientific name list use case.
type ScnmSearchUseCase interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
