package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/outbound"
)

var _ inbound.ScnmInfoUseCase = (*ScnmInfoService)(nil)

// ScnmInfoService runs the scientific name detail use case.
type ScnmInfoService struct {
	port outbound.ScnmInfoPort
}

// NewScnmInfoService creates a scientific name detail service.
func NewScnmInfoService(port outbound.ScnmInfoPort) *ScnmInfoService {
	return &ScnmInfoService{port: port}
}

// ScnmInfo returns scientific name detail information.
func (s *ScnmInfoService) ScnmInfo(ctx context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	if strings.TrimSpace(query.ReqPlantScnmID) == "" {
		return application.ScnmInfoResult{}, errors.New("reqPlantScnmId is required")
	}

	return s.port.ScnmInfo(ctx, query)
}
