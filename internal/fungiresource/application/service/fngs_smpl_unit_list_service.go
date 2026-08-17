package service

import (
	"context"
	"errors"
	"strings"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

var _ inbound.FngsSmplUnitListUseCase = (*FngsSmplUnitListService)(nil)

// FngsSmplUnitListService runs the fungi specimen detail list use case.
type FngsSmplUnitListService struct {
	port outbound.FngsSmplUnitListPort
}

// NewFngsSmplUnitListService creates a fungi specimen detail list service.
func NewFngsSmplUnitListService(port outbound.FngsSmplUnitListPort) *FngsSmplUnitListService {
	return &FngsSmplUnitListService{port: port}
}

// FngsSmplUnitList returns fungi specimen details.
func (s *FngsSmplUnitListService) FngsSmplUnitList(ctx context.Context, query application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error) {
	if query.PageNo < 1 {
		return application.FngsSmplUnitListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.FngsSmplUnitListResult{}, errors.New("numOfRows must be greater than zero")
	}
	if strings.TrimSpace(query.ReqFngsID) == "" {
		return application.FngsSmplUnitListResult{}, errors.New("reqFngsId is required")
	}

	return s.port.FngsSmplUnitList(ctx, query)
}
