package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
	"github.com/49EHyeon42/KNA-MCP/internal/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/application/port/outbound"
)

var _ inbound.PlantPictorialBookInformationUseCase = (*PlantPictorialBookInformationService)(nil)

// PlantPictorialBookInformationService runs the plant pictorial book information use case.
type PlantPictorialBookInformationService struct {
	port outbound.PlantPictorialBookInformationPort
}

// NewPlantPictorialBookInformationService creates a plant pictorial book information service.
func NewPlantPictorialBookInformationService(port outbound.PlantPictorialBookInformationPort) *PlantPictorialBookInformationService {
	return &PlantPictorialBookInformationService{port: port}
}

// PlantPictorialBookInformation returns plant pictorial book information.
func (s *PlantPictorialBookInformationService) PlantPictorialBookInformation(ctx context.Context, query application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error) {
	if strings.TrimSpace(query.RequestPlantPictorialBookNumber) == "" {
		return application.PlantPictorialBookInformationResult{}, errors.New("requestPlantPictorialBookNumber is required")
	}

	return s.port.PlantPictorialBookInformation(ctx, query)
}
