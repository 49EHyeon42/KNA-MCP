package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPilbkInfoUseCase defines the plant pictorial book information use case.
type PlantPilbkInfoUseCase interface {
	PlantPilbkInfo(context.Context, application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error)
}
