package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
)

// PlantMstnsListUseCase defines the plant miniature list use case.
type PlantMstnsListUseCase interface {
	PlantMstnsList(context.Context, application.PlantMstnsListQuery) (application.PlantMstnsListResult, error)
}
