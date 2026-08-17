package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

var _ inbound.InsectPrtctListUseCase = (*InsectPrtctListService)(nil)

// InsectPrtctListService runs the endangered insect list use case.
type InsectPrtctListService struct {
	port outbound.InsectPrtctListPort
}

// NewInsectPrtctListService creates an endangered insect list service.
func NewInsectPrtctListService(port outbound.InsectPrtctListPort) *InsectPrtctListService {
	return &InsectPrtctListService{port: port}
}

// InsectPrtctList returns endangered insects.
func (s *InsectPrtctListService) InsectPrtctList(ctx context.Context, query application.InsectPrtctListQuery) (application.InsectPrtctListResult, error) {
	if query.PageNo < 1 {
		return application.InsectPrtctListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.InsectPrtctListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.InsectPrtctList(ctx, query)
}
