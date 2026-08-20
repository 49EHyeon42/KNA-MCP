package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
)

// ScnmSearchUseCase defines the fungi scientific name list use case.
type ScnmSearchUseCase interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
