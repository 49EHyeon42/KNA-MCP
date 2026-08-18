package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

var _ inbound.FngsSmplSearchUseCase = (*FngsSmplSearchService)(nil)

// FngsSmplSearchService runs the fungi sample search use case.
type FngsSmplSearchService struct {
	port outbound.FngsSmplSearchPort
}

// NewFngsSmplSearchService creates a fungi sample search service.
func NewFngsSmplSearchService(port outbound.FngsSmplSearchPort) *FngsSmplSearchService {
	return &FngsSmplSearchService{port: port}
}

// FngsSmplSearch searches fungi samples.
func (s *FngsSmplSearchService) FngsSmplSearch(ctx context.Context, query application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error) {
	if query.PageNo < 1 {
		return application.FngsSmplSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.FngsSmplSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.FngsSmplSearch(ctx, query)
}
