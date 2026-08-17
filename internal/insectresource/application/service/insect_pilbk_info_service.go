package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

var _ inbound.InsectPilbkInfoUseCase = (*InsectPilbkInfoService)(nil)

// InsectPilbkInfoService runs the insect pictorial book detail use case.
type InsectPilbkInfoService struct {
	port outbound.InsectPilbkInfoPort
}

// NewInsectPilbkInfoService creates an insect pictorial book detail service.
func NewInsectPilbkInfoService(port outbound.InsectPilbkInfoPort) *InsectPilbkInfoService {
	return &InsectPilbkInfoService{port: port}
}

// InsectPilbkInfo returns insect pictorial book detail information.
func (s *InsectPilbkInfoService) InsectPilbkInfo(ctx context.Context, query application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error) {
	if strings.TrimSpace(query.ReqInsctPilbkNo) == "" {
		return application.InsectPilbkInfoResult{}, errors.New("reqInsctPilbkNo is required")
	}

	return s.port.InsectPilbkInfo(ctx, query)
}
