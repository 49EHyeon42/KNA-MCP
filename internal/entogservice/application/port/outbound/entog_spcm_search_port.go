package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogSpcmSearchPort defines the outbound port for entognath specimen searches.
type EntogSpcmSearchPort interface {
	EntogSpcmSearch(context.Context, application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error)
}
