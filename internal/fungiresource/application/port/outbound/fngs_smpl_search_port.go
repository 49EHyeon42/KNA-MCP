package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsSmplSearchPort defines the outbound port for fungi sample searches.
type FngsSmplSearchPort interface {
	FngsSmplSearch(context.Context, application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error)
}
