package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

var _ inbound.EntogSpcmInfoUseCase = (*EntogSpcmInfoService)(nil)

// EntogSpcmInfoService runs the entognath specimen detail lookup use case.
type EntogSpcmInfoService struct {
	port outbound.EntogSpcmInfoPort
}

// NewEntogSpcmInfoService creates an entognath specimen detail lookup service.
func NewEntogSpcmInfoService(port outbound.EntogSpcmInfoPort) *EntogSpcmInfoService {
	return &EntogSpcmInfoService{port: port}
}

// EntogSpcmInfo gets entognath specimen detail information.
func (s *EntogSpcmInfoService) EntogSpcmInfo(ctx context.Context, query application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error) {
	if strings.TrimSpace(query.Q1) == "" {
		return application.EntogSpcmInfoResult{}, errors.New("q1 is required")
	}

	return s.port.EntogSpcmInfo(ctx, query)
}
