package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/outbound"
)

var _ inbound.ChildPilbkSearchUseCase = (*ChildPilbkSearchService)(nil)

// ChildPilbkSearchService runs the child pictorial book search use case.
type ChildPilbkSearchService struct {
	port outbound.ChildPilbkSearchPort
}

// NewChildPilbkSearchService creates a child pictorial book search service.
func NewChildPilbkSearchService(port outbound.ChildPilbkSearchPort) *ChildPilbkSearchService {
	return &ChildPilbkSearchService{port: port}
}

// ChildPilbkSearch searches the child pictorial book.
func (s *ChildPilbkSearchService) ChildPilbkSearch(ctx context.Context, query application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error) {
	if query.PageNo < 1 {
		return application.ChildPilbkSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.ChildPilbkSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.ChildPilbkSearch(ctx, query)
}
