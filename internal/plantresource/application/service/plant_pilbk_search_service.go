package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantPilbkSearchUseCase = (*PlantPilbkSearchService)(nil)

// PlantPilbkSearchService runs the plant pictorial book search use case.
type PlantPilbkSearchService struct {
	port outbound.PlantPilbkSearchPort
}

// NewPlantPilbkSearchService creates a plant pictorial book search service.
func NewPlantPilbkSearchService(port outbound.PlantPilbkSearchPort) *PlantPilbkSearchService {
	return &PlantPilbkSearchService{port: port}
}

// PlantPilbkSearch searches the plant pictorial book.
func (s *PlantPilbkSearchService) PlantPilbkSearch(ctx context.Context, query application.PlantPilbkSearchQuery) (application.PlantPilbkSearchResult, error) {
	if query.PageNo < 1 {
		return application.PlantPilbkSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantPilbkSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantPilbkSearch(ctx, query)
}
