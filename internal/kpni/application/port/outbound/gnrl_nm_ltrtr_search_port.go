package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

// GnrlNmLtrtrSearchPort defines the outbound port for plant general name literature lists.
type GnrlNmLtrtrSearchPort interface {
	GnrlNmLtrtrSearch(context.Context, application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error)
}
