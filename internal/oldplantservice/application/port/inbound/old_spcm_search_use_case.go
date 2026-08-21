package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
)

// OldSpcmSearchUseCase defines the old plant specimen search use case.
type OldSpcmSearchUseCase interface {
	OldSpcmSearch(context.Context, application.OldSpcmSearchQuery) (application.OldSpcmSearchResult, error)
}
