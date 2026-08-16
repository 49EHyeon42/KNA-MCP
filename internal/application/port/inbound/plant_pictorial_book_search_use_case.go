package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

// PlantPictorialBookSearchUseCase defines the plant pictorial book search use case.
type PlantPictorialBookSearchUseCase interface {
	PlantPictorialBookSearch(context.Context, application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error)
}
