package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantFolkAreaListUseCase defines the folk plant area list use case.
type PlantFolkAreaListUseCase interface {
	PlantFolkAreaList(context.Context, application.PlantFolkAreaListQuery) (application.PlantFolkAreaListResult, error)
}
