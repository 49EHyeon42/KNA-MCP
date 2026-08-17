package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

var _ inbound.InsectSmplUnitListUseCase = (*InsectSmplUnitListService)(nil)

// InsectSmplUnitListService runs the insect specimen detail list use case.
type InsectSmplUnitListService struct {
	port outbound.InsectSmplUnitListPort
}

// NewInsectSmplUnitListService creates an insect specimen detail list service.
func NewInsectSmplUnitListService(port outbound.InsectSmplUnitListPort) *InsectSmplUnitListService {
	return &InsectSmplUnitListService{port: port}
}

// InsectSmplUnitList returns insect specimen details.
func (s *InsectSmplUnitListService) InsectSmplUnitList(ctx context.Context, query application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error) {
	if query.PageNo < 1 {
		return application.InsectSmplUnitListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.InsectSmplUnitListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.ReqInsctSpecsID) == "" {
		return application.InsectSmplUnitListResult{}, errors.New("reqInsctSpecsId is required")
	}

	return s.port.InsectSmplUnitList(ctx, query)
}
