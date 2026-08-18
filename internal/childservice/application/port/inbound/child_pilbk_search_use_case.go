package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

// ChildPilbkSearchUseCase defines the child pictorial book search use case.
type ChildPilbkSearchUseCase interface {
	ChildPilbkSearch(context.Context, application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error)
}
