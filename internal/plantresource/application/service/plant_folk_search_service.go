package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantFolkSearchUseCase = (*PlantFolkSearchService)(nil)

// PlantFolkSearchService runs the folk plant search use case.
type PlantFolkSearchService struct {
	port outbound.PlantFolkSearchPort
}

// NewPlantFolkSearchService creates a folk plant search service.
func NewPlantFolkSearchService(port outbound.PlantFolkSearchPort) *PlantFolkSearchService {
	return &PlantFolkSearchService{port: port}
}

// PlantFolkSearch searches folk plants.
func (s *PlantFolkSearchService) PlantFolkSearch(ctx context.Context, query application.PlantFolkSearchQuery) (application.PlantFolkSearchResult, error) {
	if query.PageNo < 1 {
		return application.PlantFolkSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantFolkSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantFolkSearch(ctx, query)
}
