package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

// ChildPilbkSearchPort defines the outbound port for child pictorial book searches.
type ChildPilbkSearchPort interface {
	ChildPilbkSearch(context.Context, application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error)
}
