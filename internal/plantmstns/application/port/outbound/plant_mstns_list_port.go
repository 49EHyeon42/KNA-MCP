package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
)

// PlantMstnsListPort defines the outbound port for plant miniature lists.
type PlantMstnsListPort interface {
	PlantMstnsList(context.Context, application.PlantMstnsListQuery) (application.PlantMstnsListResult, error)
}
