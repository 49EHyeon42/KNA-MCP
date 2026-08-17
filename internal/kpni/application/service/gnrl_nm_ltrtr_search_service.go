package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/outbound"
)

var _ inbound.GnrlNmLtrtrSearchUseCase = (*GnrlNmLtrtrSearchService)(nil)

// GnrlNmLtrtrSearchService runs the plant general name literature list use case.
type GnrlNmLtrtrSearchService struct {
	port outbound.GnrlNmLtrtrSearchPort
}

// NewGnrlNmLtrtrSearchService creates a plant general name literature list service.
func NewGnrlNmLtrtrSearchService(port outbound.GnrlNmLtrtrSearchPort) *GnrlNmLtrtrSearchService {
	return &GnrlNmLtrtrSearchService{port: port}
}

// GnrlNmLtrtrSearch returns plant general name literature information.
func (s *GnrlNmLtrtrSearchService) GnrlNmLtrtrSearch(ctx context.Context, query application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error) {
	if query.PageNo < 1 {
		return application.GnrlNmLtrtrSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.GnrlNmLtrtrSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.GnrlNmLtrtrSearch(ctx, query)
}
