package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsSmplSearchUseCase defines the fungi sample search use case.
type FngsSmplSearchUseCase interface {
	FngsSmplSearch(context.Context, application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error)
}
