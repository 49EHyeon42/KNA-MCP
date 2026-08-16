package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantSeedGrmntListUseCase = (*PlantSeedGrmntListService)(nil)

// PlantSeedGrmntListService runs the plant seed germination list use case.
type PlantSeedGrmntListService struct {
	port outbound.PlantSeedGrmntListPort
}

// NewPlantSeedGrmntListService creates a plant seed germination list service.
func NewPlantSeedGrmntListService(port outbound.PlantSeedGrmntListPort) *PlantSeedGrmntListService {
	return &PlantSeedGrmntListService{port: port}
}

// PlantSeedGrmntList returns plant seed germination information.
func (s *PlantSeedGrmntListService) PlantSeedGrmntList(ctx context.Context, query application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error) {
	if query.PageNo < 1 {
		return application.PlantSeedGrmntListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.PlantSeedGrmntListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.ReqSeedSpecsID) == "" {
		return application.PlantSeedGrmntListResult{}, errors.New("reqSeedSpecsId is required")
	}

	return s.port.PlantSeedGrmntList(ctx, query)
}
