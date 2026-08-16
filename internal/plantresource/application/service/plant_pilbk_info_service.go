package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

var _ inbound.PlantPilbkInfoUseCase = (*PlantPilbkInfoService)(nil)

// PlantPilbkInfoService runs the plant pictorial book information use case.
type PlantPilbkInfoService struct {
	port outbound.PlantPilbkInfoPort
}

// NewPlantPilbkInfoService creates a plant pictorial book information service.
func NewPlantPilbkInfoService(port outbound.PlantPilbkInfoPort) *PlantPilbkInfoService {
	return &PlantPilbkInfoService{port: port}
}

// PlantPilbkInfo returns plant pictorial book information.
func (s *PlantPilbkInfoService) PlantPilbkInfo(ctx context.Context, query application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error) {
	if strings.TrimSpace(query.ReqPlantPilbkNo) == "" {
		return application.PlantPilbkInfoResult{}, errors.New("reqPlantPilbkNo is required")
	}

	return s.port.PlantPilbkInfo(ctx, query)
}
