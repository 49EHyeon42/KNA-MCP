package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/port/outbound"
)

var _ inbound.ScnmInfoUseCase = (*ScnmInfoService)(nil)

// ScnmInfoService runs the fungi scientific name detail use case.
type ScnmInfoService struct {
	port outbound.ScnmInfoPort
}

// NewScnmInfoService creates a fungi scientific name detail service.
func NewScnmInfoService(port outbound.ScnmInfoPort) *ScnmInfoService {
	return &ScnmInfoService{port: port}
}

// ScnmInfo returns fungi scientific name detail information.
func (s *ScnmInfoService) ScnmInfo(ctx context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	if strings.TrimSpace(query.ReqFngsScnmID) == "" {
		return application.ScnmInfoResult{}, errors.New("reqFngsScnmId is required")
	}

	return s.port.ScnmInfo(ctx, query)
}
