package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/application/port/outbound"
)

var _ inbound.ScnmInfoUseCase = (*ScnmInfoService)(nil)

// ScnmInfoService runs the insect scientific name detail use case.
type ScnmInfoService struct {
	port outbound.ScnmInfoPort
}

// NewScnmInfoService creates an insect scientific name detail service.
func NewScnmInfoService(port outbound.ScnmInfoPort) *ScnmInfoService {
	return &ScnmInfoService{port: port}
}

// ScnmInfo returns insect scientific name detail information.
func (s *ScnmInfoService) ScnmInfo(ctx context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	if strings.TrimSpace(query.ReqInsctScnmID) == "" {
		return application.ScnmInfoResult{}, errors.New("reqInsctScnmId is required")
	}

	return s.port.ScnmInfo(ctx, query)
}
