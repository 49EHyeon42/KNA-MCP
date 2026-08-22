package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

var _ inbound.EntogIlstrSearchUseCase = (*EntogIlstrSearchService)(nil)

// EntogIlstrSearchService runs the entognath pictorial book search use case.
type EntogIlstrSearchService struct {
	port outbound.EntogIlstrSearchPort
}

// NewEntogIlstrSearchService creates an entognath pictorial book search service.
func NewEntogIlstrSearchService(port outbound.EntogIlstrSearchPort) *EntogIlstrSearchService {
	return &EntogIlstrSearchService{port: port}
}

// EntogIlstrSearch searches the entognath pictorial book.
func (s *EntogIlstrSearchService) EntogIlstrSearch(ctx context.Context, query application.EntogIlstrSearchQuery) (application.EntogIlstrSearchResult, error) {
	if query.St != "1" && query.St != "2" && query.St != "3" && query.St != "4" {
		return application.EntogIlstrSearchResult{}, errors.New("st must be one of 1, 2, 3, or 4")
	}
	if query.Sw == "" {
		return application.EntogIlstrSearchResult{}, errors.New("sw is required")
	}
	if query.PageNo < 1 {
		return application.EntogIlstrSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.EntogIlstrSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.EntogIlstrSearch(ctx, query)
}
