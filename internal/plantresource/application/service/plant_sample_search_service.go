package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSampleSearchUseCase = (*PlantSampleSearchService)(nil)

// PlantSampleSearchService runs the plant sample search use case.
type PlantSampleSearchService struct {
	port outbound.PlantSampleSearchPort
}

// NewPlantSampleSearchService creates a plant sample search service.
func NewPlantSampleSearchService(port outbound.PlantSampleSearchPort) *PlantSampleSearchService {
	return &PlantSampleSearchService{port: port}
}

// PlantSampleSearch searches plant samples.
func (s *PlantSampleSearchService) PlantSampleSearch(ctx context.Context, query application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error) {
	if query.PageNumber < 1 {
		return application.PlantSampleSearchResult{}, errors.New("pageNumber must be greater than zero")
	}
	if query.NumberOfRows < 1 {
		return application.PlantSampleSearchResult{}, errors.New("numberOfRows must be greater than zero")
	}

	return s.port.PlantSampleSearch(ctx, query)
}
