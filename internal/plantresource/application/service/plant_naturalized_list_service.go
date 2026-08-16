package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantNaturalizedListUseCase = (*PlantNaturalizedListService)(nil)

// PlantNaturalizedListService runs the naturalized plant list use case.
type PlantNaturalizedListService struct {
	port outbound.PlantNaturalizedListPort
}

// NewPlantNaturalizedListService creates a naturalized plant list service.
func NewPlantNaturalizedListService(port outbound.PlantNaturalizedListPort) *PlantNaturalizedListService {
	return &PlantNaturalizedListService{port: port}
}

// PlantNaturalizedList returns naturalized plant information.
func (s *PlantNaturalizedListService) PlantNaturalizedList(ctx context.Context, query application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error) {
	if query.PageNo < 1 {
		return application.PlantNaturalizedListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantNaturalizedListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantNaturalizedList(ctx, query)
}
