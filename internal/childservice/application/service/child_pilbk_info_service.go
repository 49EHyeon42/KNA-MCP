package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/outbound"
)

var _ inbound.ChildPilbkInfoUseCase = (*ChildPilbkInfoService)(nil)

// ChildPilbkInfoService runs the child pictorial book detail use case.
type ChildPilbkInfoService struct {
	port outbound.ChildPilbkInfoPort
}

// NewChildPilbkInfoService creates a child pictorial book detail service.
func NewChildPilbkInfoService(port outbound.ChildPilbkInfoPort) *ChildPilbkInfoService {
	return &ChildPilbkInfoService{port: port}
}

// ChildPilbkInfo returns child pictorial book detail information.
func (s *ChildPilbkInfoService) ChildPilbkInfo(ctx context.Context, query application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error) {
	if strings.TrimSpace(query.ReqChildLvbngPilbkNo) == "" {
		return application.ChildPilbkInfoResult{}, errors.New("reqChildLvbngPilbkNo is required")
	}

	return s.port.ChildPilbkInfo(ctx, query)
}
