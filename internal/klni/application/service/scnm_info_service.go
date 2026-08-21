package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/outbound"
)

var _ inbound.ScnmInfoUseCase = (*ScnmInfoService)(nil)

// ScnmInfoService runs the lichen scientific name detail use case.
type ScnmInfoService struct {
	port outbound.ScnmInfoPort
}

// NewScnmInfoService creates a lichen scientific name detail service.
func NewScnmInfoService(port outbound.ScnmInfoPort) *ScnmInfoService {
	return &ScnmInfoService{port: port}
}

// ScnmInfo returns lichen scientific name detail information.
func (s *ScnmInfoService) ScnmInfo(ctx context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	if strings.TrimSpace(query.ReqLchnScnmID) == "" {
		return application.ScnmInfoResult{}, errors.New("reqLchnScnmId is required")
	}

	return s.port.ScnmInfo(ctx, query)
}
