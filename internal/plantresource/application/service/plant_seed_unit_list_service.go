package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSeedUnitListUseCase = (*PlantSeedUnitListService)(nil)

// PlantSeedUnitListService runs the plant seed unit list use case.
type PlantSeedUnitListService struct {
	port outbound.PlantSeedUnitListPort
}

// NewPlantSeedUnitListService creates a plant seed unit list service.
func NewPlantSeedUnitListService(port outbound.PlantSeedUnitListPort) *PlantSeedUnitListService {
	return &PlantSeedUnitListService{port: port}
}

// PlantSeedUnitList returns plant seed unit information.
func (s *PlantSeedUnitListService) PlantSeedUnitList(ctx context.Context, query application.PlantSeedUnitListQuery) (application.PlantSeedUnitListResult, error) {
	if query.PageNo < 1 {
		return application.PlantSeedUnitListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSeedUnitListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.ReqSeedSpecsID) == "" {
		return application.PlantSeedUnitListResult{}, errors.New("reqSeedSpecsId is required")
	}

	return s.port.PlantSeedUnitList(ctx, query)
}
