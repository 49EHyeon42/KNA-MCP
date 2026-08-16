package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantWordListPort defines the outbound port for plant word lists.
type PlantWordListPort interface {
	PlantWordList(context.Context, application.PlantWordListQuery) (application.PlantWordListResult, error)
}
