package inbound

import (
	"context"

	"kna-mcp/internal/application"
)

// PlantPictorialBookSearchUseCase defines the plant pictorial book search use case.
type PlantPictorialBookSearchUseCase interface {
	PlantPictorialBookSearch(context.Context, application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error)
}
