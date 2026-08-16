package service

import (
	"context"
	"errors"

	"kna-mcp/internal/application"
	"kna-mcp/internal/application/port/inbound"
	"kna-mcp/internal/application/port/outbound"
)

var _ inbound.PlantPictorialBookSearchUseCase = (*PlantPictorialBookSearchService)(nil)

// PlantPictorialBookSearchService runs the plant pictorial book search use case.
type PlantPictorialBookSearchService struct {
	port outbound.PlantPictorialBookSearchPort
}

// NewPlantPictorialBookSearchService creates a plant pictorial book search service.
func NewPlantPictorialBookSearchService(port outbound.PlantPictorialBookSearchPort) *PlantPictorialBookSearchService {
	return &PlantPictorialBookSearchService{port: port}
}

// PlantPictorialBookSearch searches the plant pictorial book.
func (s *PlantPictorialBookSearchService) PlantPictorialBookSearch(ctx context.Context, query application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error) {
	if query.PageNumber < 1 {
		return application.PlantPictorialBookSearchResult{}, errors.New("pageNumber must be greater than zero")
	}
	if query.NumberOfRows < 1 {
		return application.PlantPictorialBookSearchResult{}, errors.New("numberOfRows must be greater than zero")
	}

	return s.port.PlantPictorialBookSearch(ctx, query)
}
