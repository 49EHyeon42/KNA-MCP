package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSmplSearchUseCase = (*PlantSmplSearchService)(nil)

// PlantSmplSearchService runs the plant sample search use case.
type PlantSmplSearchService struct {
	port outbound.PlantSmplSearchPort
}

// NewPlantSmplSearchService creates a plant sample search service.
func NewPlantSmplSearchService(port outbound.PlantSmplSearchPort) *PlantSmplSearchService {
	return &PlantSmplSearchService{port: port}
}

// PlantSmplSearch searches plant samples.
func (s *PlantSmplSearchService) PlantSmplSearch(ctx context.Context, query application.PlantSmplSearchQuery) (application.PlantSmplSearchResult, error) {
	if query.PageNo < 1 {
		return application.PlantSmplSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSmplSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantSmplSearch(ctx, query)
}
