package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

var _ inbound.FngsPilbkSearchUseCase = (*FngsPilbkSearchService)(nil)

// FngsPilbkSearchService runs the fungi pictorial book search use case.
type FngsPilbkSearchService struct {
	port outbound.FngsPilbkSearchPort
}

// NewFngsPilbkSearchService creates a fungi pictorial book search service.
func NewFngsPilbkSearchService(port outbound.FngsPilbkSearchPort) *FngsPilbkSearchService {
	return &FngsPilbkSearchService{port: port}
}

// FngsPilbkSearch searches the fungi pictorial book.
func (s *FngsPilbkSearchService) FngsPilbkSearch(ctx context.Context, query application.FngsPilbkSearchQuery) (application.FngsPilbkSearchResult, error) {
	if query.PageNo < 1 {
		return application.FngsPilbkSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.FngsPilbkSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.FngsPilbkSearch(ctx, query)
}
