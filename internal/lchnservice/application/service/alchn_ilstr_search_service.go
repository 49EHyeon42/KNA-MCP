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

var _ inbound.AlchnIlstrSearchUseCase = (*AlchnIlstrSearchService)(nil)

// AlchnIlstrSearchService runs the lichen pictorial book search use case.
type AlchnIlstrSearchService struct {
	port outbound.AlchnIlstrSearchPort
}

// NewAlchnIlstrSearchService creates a lichen pictorial book search service.
func NewAlchnIlstrSearchService(port outbound.AlchnIlstrSearchPort) *AlchnIlstrSearchService {
	return &AlchnIlstrSearchService{port: port}
}

// AlchnIlstrSearch searches the lichen pictorial book.
func (s *AlchnIlstrSearchService) AlchnIlstrSearch(ctx context.Context, query application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error) {
	if query.St != "1" && query.St != "2" && query.St != "3" && query.St != "4" {
		return application.AlchnIlstrSearchResult{}, errors.New("st must be one of 1, 2, 3, or 4")
	}
	if query.Sw == "" {
		return application.AlchnIlstrSearchResult{}, errors.New("sw is required")
	}
	if query.PageNo < 1 {
		return application.AlchnIlstrSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.AlchnIlstrSearchResult{}, errors.New("numOfRows must be greater than zero")
	}
	if query.DateGbn == "" {
		if query.DateFrom != "" || query.DateTo != "" {
			return application.AlchnIlstrSearchResult{}, errors.New("dateGbn is required when dateFrom or dateTo is provided")
		}
		return s.port.AlchnIlstrSearch(ctx, query)
	}
	if query.DateGbn != "1" && query.DateGbn != "2" {
		return application.AlchnIlstrSearchResult{}, errors.New("dateGbn must be 1 or 2")
	}
	dateFrom, err := parseAlchnIlstrSearchDate("dateFrom", query.DateFrom)
	if err != nil {
		return application.AlchnIlstrSearchResult{}, err
	}
	dateTo, err := parseAlchnIlstrSearchDate("dateTo", query.DateTo)
	if err != nil {
		return application.AlchnIlstrSearchResult{}, err
	}
	if dateFrom.After(dateTo) {
		return application.AlchnIlstrSearchResult{}, errors.New("dateFrom must not be after dateTo")
	}

	return s.port.AlchnIlstrSearch(ctx, query)
}

func parseAlchnIlstrSearchDate(name, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required when dateGbn is provided", name)
	}
	date, err := time.Parse("20060102", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be a valid date in yyyyMMdd format", name)
	}
	return date, nil
}
