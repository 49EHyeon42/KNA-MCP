package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
)

// ScnmSearchUseCase defines the insect scientific name list use case.
type ScnmSearchUseCase interface {
	ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error)
}
