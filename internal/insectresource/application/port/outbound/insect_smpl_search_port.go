package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectSmplSearchPort defines the outbound port for insect sample searches.
type InsectSmplSearchPort interface {
	InsectSmplSearch(context.Context, application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error)
}
