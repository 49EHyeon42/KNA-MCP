package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnIlstrSearchUseCase defines the lichen pictorial book search use case.
type AlchnIlstrSearchUseCase interface {
	AlchnIlstrSearch(context.Context, application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error)
}
