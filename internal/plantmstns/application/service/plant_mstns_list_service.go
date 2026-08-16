package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/outbound"
)

var _ inbound.PlantMstnsListUseCase = (*PlantMstnsListService)(nil)

// PlantMstnsListService runs the plant miniature list use case.
type PlantMstnsListService struct {
	port outbound.PlantMstnsListPort
}

// NewPlantMstnsListService creates a plant miniature list service.
func NewPlantMstnsListService(port outbound.PlantMstnsListPort) *PlantMstnsListService {
	return &PlantMstnsListService{port: port}
}

// PlantMstnsList returns plant miniature information.
func (s *PlantMstnsListService) PlantMstnsList(ctx context.Context, query application.PlantMstnsListQuery) (application.PlantMstnsListResult, error) {
	if query.PageNo < 1 {
		return application.PlantMstnsListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantMstnsListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantMstnsList(ctx, query)
}
