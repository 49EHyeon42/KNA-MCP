package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantFolkAreaListUseCase = (*PlantFolkAreaListService)(nil)

// PlantFolkAreaListService runs the folk plant area list use case.
type PlantFolkAreaListService struct {
	port outbound.PlantFolkAreaListPort
}

// NewPlantFolkAreaListService creates a folk plant area list service.
func NewPlantFolkAreaListService(port outbound.PlantFolkAreaListPort) *PlantFolkAreaListService {
	return &PlantFolkAreaListService{port: port}
}

// PlantFolkAreaList returns folk plant area information.
func (s *PlantFolkAreaListService) PlantFolkAreaList(ctx context.Context, query application.PlantFolkAreaListQuery) (application.PlantFolkAreaListResult, error) {
	if query.PageNo < 1 {
		return application.PlantFolkAreaListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantFolkAreaListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.FlpltID) == "" {
		return application.PlantFolkAreaListResult{}, errors.New("flpltId is required")
	}

	return s.port.PlantFolkAreaList(ctx, query)
}
