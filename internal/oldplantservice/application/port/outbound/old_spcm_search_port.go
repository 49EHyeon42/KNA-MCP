package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
)

// OldSpcmSearchPort defines the outbound port for old plant specimen searches.
type OldSpcmSearchPort interface {
	OldSpcmSearch(context.Context, application.OldSpcmSearchQuery) (application.OldSpcmSearchResult, error)
}
