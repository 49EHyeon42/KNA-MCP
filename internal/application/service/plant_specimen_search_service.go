package service

import (
	"context"
	"errors"

	"kna-mcp/internal/application"
	"kna-mcp/internal/application/port/inbound"
	"kna-mcp/internal/application/port/outbound"
)

var _ inbound.PlantSpecimenSearchUseCase = (*PlantSpecimenSearchService)(nil)

// PlantSpecimenSearchService runs the plant specimen search use case.
type PlantSpecimenSearchService struct {
	port outbound.PlantSpecimenSearchPort
}

// NewPlantSpecimenSearchService creates a plant specimen search service.
func NewPlantSpecimenSearchService(port outbound.PlantSpecimenSearchPort) *PlantSpecimenSearchService {
	return &PlantSpecimenSearchService{port: port}
}

// PlantSpecimenSearch searches plant specimens.
func (s *PlantSpecimenSearchService) PlantSpecimenSearch(ctx context.Context, query application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error) {
	if query.PageNumber < 1 {
		return application.PlantSpecimenSearchResult{}, errors.New("pageNumber must be greater than zero")
	}
	if query.NumberOfRows < 1 {
		return application.PlantSpecimenSearchResult{}, errors.New("numberOfRows must be greater than zero")
	}

	return s.port.PlantSpecimenSearch(ctx, query)
}
