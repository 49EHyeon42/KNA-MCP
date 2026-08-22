package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogIlstrSearchUseCase defines the entognath pictorial book search use case.
type EntogIlstrSearchUseCase interface {
	EntogIlstrSearch(context.Context, application.EntogIlstrSearchQuery) (application.EntogIlstrSearchResult, error)
}
