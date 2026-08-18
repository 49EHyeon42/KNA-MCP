package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsPilbkSearchUseCase defines the fungi pictorial book search use case.
type FngsPilbkSearchUseCase interface {
	FngsPilbkSearch(context.Context, application.FngsPilbkSearchQuery) (application.FngsPilbkSearchResult, error)
}
