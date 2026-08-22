package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogSpcmSearchUseCase defines the inbound port for entognath specimen searches.
type EntogSpcmSearchUseCase interface {
	EntogSpcmSearch(context.Context, application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error)
}
