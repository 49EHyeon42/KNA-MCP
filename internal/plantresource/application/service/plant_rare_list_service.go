package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantRareListUseCase = (*PlantRareListService)(nil)

// PlantRareListService runs the rare plant list use case.
type PlantRareListService struct {
	port outbound.PlantRareListPort
}

// NewPlantRareListService creates a rare plant list service.
func NewPlantRareListService(port outbound.PlantRareListPort) *PlantRareListService {
	return &PlantRareListService{port: port}
}

// PlantRareList returns rare plant information.
func (s *PlantRareListService) PlantRareList(ctx context.Context, query application.PlantRareListQuery) (application.PlantRareListResult, error) {
	if query.PageNo < 1 {
		return application.PlantRareListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantRareListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantRareList(ctx, query)
}
