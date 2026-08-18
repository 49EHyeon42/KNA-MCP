package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnSpcmSearchUseCase defines the lichen specimen search use case.
type AlchnSpcmSearchUseCase interface {
	AlchnSpcmSearch(context.Context, application.AlchnSpcmSearchQuery) (application.AlchnSpcmSearchResult, error)
}
