package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

// GnrlNmLtrtrSearchUseCase defines the plant general name literature list use case.
type GnrlNmLtrtrSearchUseCase interface {
	GnrlNmLtrtrSearch(context.Context, application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error)
}
