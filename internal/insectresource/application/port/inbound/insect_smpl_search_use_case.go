package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectSmplSearchUseCase defines the insect sample search use case.
type InsectSmplSearchUseCase interface {
	InsectSmplSearch(context.Context, application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error)
}
