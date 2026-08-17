package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPilbkSearchPort defines the outbound port for insect pictorial book searches.
type InsectPilbkSearchPort interface {
	InsectPilbkSearch(context.Context, application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error)
}
