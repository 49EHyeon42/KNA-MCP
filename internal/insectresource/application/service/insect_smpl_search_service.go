package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

var _ inbound.InsectSmplSearchUseCase = (*InsectSmplSearchService)(nil)

// InsectSmplSearchService runs the insect sample search use case.
type InsectSmplSearchService struct {
	port outbound.InsectSmplSearchPort
}

// NewInsectSmplSearchService creates an insect sample search service.
func NewInsectSmplSearchService(port outbound.InsectSmplSearchPort) *InsectSmplSearchService {
	return &InsectSmplSearchService{port: port}
}

// InsectSmplSearch searches insect samples.
func (s *InsectSmplSearchService) InsectSmplSearch(ctx context.Context, query application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error) {
	if query.PageNo < 1 {
		return application.InsectSmplSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.InsectSmplSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.InsectSmplSearch(ctx, query)
}
