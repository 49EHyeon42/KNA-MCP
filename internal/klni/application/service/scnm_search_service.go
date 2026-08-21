package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/outbound"
)

var _ inbound.ScnmSearchUseCase = (*ScnmSearchService)(nil)

// ScnmSearchService runs the lichen scientific name list use case.
type ScnmSearchService struct {
	port outbound.ScnmSearchPort
}

// NewScnmSearchService creates a lichen scientific name list service.
func NewScnmSearchService(port outbound.ScnmSearchPort) *ScnmSearchService {
	return &ScnmSearchService{port: port}
}

// ScnmSearch returns lichen scientific name information.
func (s *ScnmSearchService) ScnmSearch(ctx context.Context, query application.ScnmSearchQuery) (application.ScnmSearchResult, error) {
	if query.PageNo < 1 {
		return application.ScnmSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.ScnmSearchResult{}, errors.New("numOfRows must be greater than zero")
	}
	if err := validateScnmSearchDate("dateFrom", query.DateFrom); err != nil {
		return application.ScnmSearchResult{}, err
	}
	if err := validateScnmSearchDate("dateTo", query.DateTo); err != nil {
		return application.ScnmSearchResult{}, err
	}

	return s.port.ScnmSearch(ctx, query)
}

func validateScnmSearchDate(name, value string) error {
	if value == "" {
		return nil
	}
	if _, err := time.Parse("20060102", value); err != nil {
		return fmt.Errorf("%s must be a valid date in yyyyMMdd format", name)
	}
	return nil
}
