package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSmplUnitListUseCase = (*PlantSmplUnitListService)(nil)

// PlantSmplUnitListService runs the plant specimen detail list use case.
type PlantSmplUnitListService struct {
	port outbound.PlantSmplUnitListPort
}

// NewPlantSmplUnitListService creates a plant specimen detail list service.
func NewPlantSmplUnitListService(port outbound.PlantSmplUnitListPort) *PlantSmplUnitListService {
	return &PlantSmplUnitListService{port: port}
}

// PlantSmplUnitList returns plant specimen details.
func (s *PlantSmplUnitListService) PlantSmplUnitList(ctx context.Context, query application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error) {
	if query.PageNo < 1 {
		return application.PlantSmplUnitListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSmplUnitListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.ReqPlantSpecsID) == "" {
		return application.PlantSmplUnitListResult{}, errors.New("reqPlantSpecsId is required")
	}

	return s.port.PlantSmplUnitList(ctx, query)
}
