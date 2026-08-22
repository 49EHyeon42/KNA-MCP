package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogIlstrSearchPort defines the outbound port for entognath pictorial book searches.
type EntogIlstrSearchPort interface {
	EntogIlstrSearch(context.Context, application.EntogIlstrSearchQuery) (application.EntogIlstrSearchResult, error)
}
