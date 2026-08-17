package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

var _ inbound.FngsPilbkInfoUseCase = (*FngsPilbkInfoService)(nil)

// FngsPilbkInfoService runs the fungi pictorial book detail use case.
type FngsPilbkInfoService struct {
	port outbound.FngsPilbkInfoPort
}

// NewFngsPilbkInfoService creates a fungi pictorial book detail service.
func NewFngsPilbkInfoService(port outbound.FngsPilbkInfoPort) *FngsPilbkInfoService {
	return &FngsPilbkInfoService{port: port}
}

// FngsPilbkInfo returns fungi pictorial book detail information.
func (s *FngsPilbkInfoService) FngsPilbkInfo(ctx context.Context, query application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error) {
	if strings.TrimSpace(query.ReqFngsPilbkNo) == "" {
		return application.FngsPilbkInfoResult{}, errors.New("reqFngsPilbkNo is required")
	}

	return s.port.FngsPilbkInfo(ctx, query)
}
