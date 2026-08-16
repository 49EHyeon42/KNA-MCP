package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSeedSearchUseCase = (*PlantSeedSearchService)(nil)

// PlantSeedSearchService runs the plant seed search use case.
type PlantSeedSearchService struct {
	port outbound.PlantSeedSearchPort
}

// NewPlantSeedSearchService creates a plant seed search service.
func NewPlantSeedSearchService(port outbound.PlantSeedSearchPort) *PlantSeedSearchService {
	return &PlantSeedSearchService{port: port}
}

// PlantSeedSearch searches plant seed information.
func (s *PlantSeedSearchService) PlantSeedSearch(ctx context.Context, query application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error) {
	if query.PageNo < 1 {
		return application.PlantSeedSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSeedSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.PlantSeedSearch(ctx, query)
}
