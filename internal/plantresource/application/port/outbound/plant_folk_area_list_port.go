package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantFolkAreaListPort defines the outbound port for folk plant area lists.
type PlantFolkAreaListPort interface {
	PlantFolkAreaList(context.Context, application.PlantFolkAreaListQuery) (application.PlantFolkAreaListResult, error)
}
