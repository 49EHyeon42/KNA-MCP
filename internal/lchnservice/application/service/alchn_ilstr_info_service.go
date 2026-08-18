package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/outbound"
)

var _ inbound.AlchnIlstrInfoUseCase = (*AlchnIlstrInfoService)(nil)

// AlchnIlstrInfoService runs the lichen pictorial book detail use case.
type AlchnIlstrInfoService struct {
	port outbound.AlchnIlstrInfoPort
}

// NewAlchnIlstrInfoService creates a lichen pictorial book detail service.
func NewAlchnIlstrInfoService(port outbound.AlchnIlstrInfoPort) *AlchnIlstrInfoService {
	return &AlchnIlstrInfoService{port: port}
}

// AlchnIlstrInfo returns lichen pictorial book detail information.
func (s *AlchnIlstrInfoService) AlchnIlstrInfo(ctx context.Context, query application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error) {
	if strings.TrimSpace(query.Q1) == "" {
		return application.AlchnIlstrInfoResult{}, errors.New("q1 is required")
	}

	return s.port.AlchnIlstrInfo(ctx, query)
}
