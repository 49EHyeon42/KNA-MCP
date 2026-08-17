package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPilbkSearchUseCase defines the insect pictorial book search use case.
type InsectPilbkSearchUseCase interface {
	InsectPilbkSearch(context.Context, application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error)
}
