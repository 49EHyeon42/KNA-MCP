package service

import (
	"context"
	"errors"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

var _ inbound.EntogSpcmSearchUseCase = (*EntogSpcmSearchService)(nil)

// EntogSpcmSearchService runs the entognath specimen search use case.
type EntogSpcmSearchService struct {
	port outbound.EntogSpcmSearchPort
}

// NewEntogSpcmSearchService creates an entognath specimen search service.
func NewEntogSpcmSearchService(port outbound.EntogSpcmSearchPort) *EntogSpcmSearchService {
	return &EntogSpcmSearchService{port: port}
}

// EntogSpcmSearch searches entognath specimens.
func (s *EntogSpcmSearchService) EntogSpcmSearch(ctx context.Context, query application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error) {
	if query.St != "1" && query.St != "2" && query.St != "3" && query.St != "4" {
		return application.EntogSpcmSearchResult{}, errors.New("st must be one of 1, 2, 3, or 4")
	}
	if query.Sw == "" {
		return application.EntogSpcmSearchResult{}, errors.New("sw is required")
	}
	if query.PageNo < 1 {
		return application.EntogSpcmSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.EntogSpcmSearchResult{}, errors.New("numOfRows must be greater than zero")
	}

	return s.port.EntogSpcmSearch(ctx, query)
}
