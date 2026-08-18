package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/outbound"
)

var _ inbound.AlchnSpcmSearchUseCase = (*AlchnSpcmSearchService)(nil)

// AlchnSpcmSearchService runs the lichen specimen search use case.
type AlchnSpcmSearchService struct {
	port outbound.AlchnSpcmSearchPort
}

// NewAlchnSpcmSearchService creates a lichen specimen search service.
func NewAlchnSpcmSearchService(port outbound.AlchnSpcmSearchPort) *AlchnSpcmSearchService {
	return &AlchnSpcmSearchService{port: port}
}

// AlchnSpcmSearch searches lichen specimens.
func (s *AlchnSpcmSearchService) AlchnSpcmSearch(ctx context.Context, query application.AlchnSpcmSearchQuery) (application.AlchnSpcmSearchResult, error) {
	if query.St != "1" && query.St != "2" && query.St != "3" && query.St != "4" {
		return application.AlchnSpcmSearchResult{}, errors.New("st must be one of 1, 2, 3, or 4")
	}
	if query.Sw == "" {
		return application.AlchnSpcmSearchResult{}, errors.New("sw is required")
	}
	if query.PageNo < 1 {
		return application.AlchnSpcmSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.AlchnSpcmSearchResult{}, errors.New("numOfRows must be greater than zero")
	}
	if query.DateGbn == "" {
		if query.DateFrom != "" || query.DateTo != "" {
			return application.AlchnSpcmSearchResult{}, errors.New("dateGbn is required when dateFrom or dateTo is provided")
		}
		return s.port.AlchnSpcmSearch(ctx, query)
	}
	if query.DateGbn != "1" && query.DateGbn != "2" {
		return application.AlchnSpcmSearchResult{}, errors.New("dateGbn must be 1 or 2")
	}
	dateFrom, err := parseAlchnSpcmSearchDate("dateFrom", query.DateFrom)
	if err != nil {
		return application.AlchnSpcmSearchResult{}, err
	}
	dateTo, err := parseAlchnSpcmSearchDate("dateTo", query.DateTo)
	if err != nil {
		return application.AlchnSpcmSearchResult{}, err
	}
	if dateFrom.After(dateTo) {
		return application.AlchnSpcmSearchResult{}, errors.New("dateFrom must not be after dateTo")
	}

	return s.port.AlchnSpcmSearch(ctx, query)
}

func parseAlchnSpcmSearchDate(name, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required when dateGbn is provided", name)
	}
	date, err := time.Parse("20060102", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be a valid date in yyyyMMdd format", name)
	}
	return date, nil
}
