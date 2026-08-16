package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantWordListUseCase = (*PlantWordListService)(nil)

// PlantWordListService runs the plant word list use case.
type PlantWordListService struct {
	port outbound.PlantWordListPort
}

// NewPlantWordListService creates a plant word list service.
func NewPlantWordListService(port outbound.PlantWordListPort) *PlantWordListService {
	return &PlantWordListService{port: port}
}

// PlantWordList returns plant word information.
func (s *PlantWordListService) PlantWordList(ctx context.Context, query application.PlantWordListQuery) (application.PlantWordListResult, error) {
	if query.PageNo < 1 {
		return application.PlantWordListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantWordListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantWordList(ctx, query)
}
