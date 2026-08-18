package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/outbound"
)

var _ inbound.AlchnSpcmInfoUseCase = (*AlchnSpcmInfoService)(nil)

// AlchnSpcmInfoService runs the lichen specimen detail use case.
type AlchnSpcmInfoService struct {
	port outbound.AlchnSpcmInfoPort
}

// NewAlchnSpcmInfoService creates a lichen specimen detail service.
func NewAlchnSpcmInfoService(port outbound.AlchnSpcmInfoPort) *AlchnSpcmInfoService {
	return &AlchnSpcmInfoService{port: port}
}

// AlchnSpcmInfo returns lichen specimen detail information.
func (s *AlchnSpcmInfoService) AlchnSpcmInfo(ctx context.Context, query application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error) {
	if strings.TrimSpace(query.Q1) == "" {
		return application.AlchnSpcmInfoResult{}, errors.New("q1 is required")
	}

	return s.port.AlchnSpcmInfo(ctx, query)
}
