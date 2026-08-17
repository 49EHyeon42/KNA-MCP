package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

var _ inbound.InsectPilbkSearchUseCase = (*InsectPilbkSearchService)(nil)

// InsectPilbkSearchService runs the insect pictorial book search use case.
type InsectPilbkSearchService struct {
	port outbound.InsectPilbkSearchPort
}

// NewInsectPilbkSearchService creates an insect pictorial book search service.
func NewInsectPilbkSearchService(port outbound.InsectPilbkSearchPort) *InsectPilbkSearchService {
	return &InsectPilbkSearchService{port: port}
}

// InsectPilbkSearch searches the insect pictorial book.
func (s *InsectPilbkSearchService) InsectPilbkSearch(ctx context.Context, query application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error) {
	if query.PageNo < 1 {
		return application.InsectPilbkSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.InsectPilbkSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.InsectPilbkSearch(ctx, query)
}
