package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/port/outbound"
)

var _ inbound.RelatedSiteListUseCase = (*RelatedSiteListService)(nil)

// RelatedSiteListService runs the related site list use case.
type RelatedSiteListService struct {
	port outbound.RelatedSiteListPort
}

// NewRelatedSiteListService creates a related site list service.
func NewRelatedSiteListService(port outbound.RelatedSiteListPort) *RelatedSiteListService {
	return &RelatedSiteListService{port: port}
}

// RelatedSiteList returns related site information.
func (s *RelatedSiteListService) RelatedSiteList(ctx context.Context, query application.RelatedSiteListQuery) (application.RelatedSiteListResult, error) {
	if query.PageNo < 1 {
		return application.RelatedSiteListResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.RelatedSiteListResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.RelatedSiteList(ctx, query)
}
