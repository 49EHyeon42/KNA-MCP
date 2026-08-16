package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPilbkInfoPort defines the outbound port for plant pictorial book information.
type PlantPilbkInfoPort interface {
	PlantPilbkInfo(context.Context, application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error)
}
