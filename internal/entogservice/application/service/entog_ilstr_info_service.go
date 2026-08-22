package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

var _ inbound.EntogIlstrInfoUseCase = (*EntogIlstrInfoService)(nil)

// EntogIlstrInfoService runs the entognath pictorial book detail use case.
type EntogIlstrInfoService struct {
	port outbound.EntogIlstrInfoPort
}

// NewEntogIlstrInfoService creates an entognath pictorial book detail service.
func NewEntogIlstrInfoService(port outbound.EntogIlstrInfoPort) *EntogIlstrInfoService {
	return &EntogIlstrInfoService{port: port}
}

// EntogIlstrInfo returns entognath pictorial book detail information.
func (s *EntogIlstrInfoService) EntogIlstrInfo(ctx context.Context, query application.EntogIlstrInfoQuery) (application.EntogIlstrInfoResult, error) {
	if strings.TrimSpace(query.Q1) == "" {
		return application.EntogIlstrInfoResult{}, errors.New("q1 is required")
	}

	return s.port.EntogIlstrInfo(ctx, query)
}
