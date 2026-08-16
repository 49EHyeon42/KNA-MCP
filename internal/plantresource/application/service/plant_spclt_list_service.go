package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSpcltListUseCase = (*PlantSpcltListService)(nil)

// PlantSpcltListService runs the endemic plant list use case.
type PlantSpcltListService struct {
	port outbound.PlantSpcltListPort
}

// NewPlantSpcltListService creates an endemic plant list service.
func NewPlantSpcltListService(port outbound.PlantSpcltListPort) *PlantSpcltListService {
	return &PlantSpcltListService{port: port}
}

// PlantSpcltList returns endemic plant information.
func (s *PlantSpcltListService) PlantSpcltList(ctx context.Context, query application.PlantSpcltListQuery) (application.PlantSpcltListResult, error) {
	if query.PageNo < 1 {
		return application.PlantSpcltListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSpcltListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantSpcltList(ctx, query)
}
