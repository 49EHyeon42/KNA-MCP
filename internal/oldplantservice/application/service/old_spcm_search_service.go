package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/port/inbound"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/port/outbound"
)

var _ inbound.OldSpcmSearchUseCase = (*OldSpcmSearchService)(nil)

// OldSpcmSearchService runs the old plant specimen search use case.
type OldSpcmSearchService struct {
	port outbound.OldSpcmSearchPort
}

// NewOldSpcmSearchService creates an old plant specimen search service.
func NewOldSpcmSearchService(port outbound.OldSpcmSearchPort) *OldSpcmSearchService {
	return &OldSpcmSearchService{port: port}
}

// OldSpcmSearch searches old plant specimens.
func (s *OldSpcmSearchService) OldSpcmSearch(ctx context.Context, query application.OldSpcmSearchQuery) (application.OldSpcmSearchResult, error) {
	if query.St != "1" && query.St != "2" {
		return application.OldSpcmSearchResult{}, errors.New("st must be 1 or 2")
	}
	if query.Sw == "" {
		return application.OldSpcmSearchResult{}, errors.New("sw is required")
	}
	if query.PageNo < 1 {
		return application.OldSpcmSearchResult{}, errors.New("pageNo must be greater than zero")
	}
	if query.NumOfRows < 1 {
		return application.OldSpcmSearchResult{}, errors.New("numOfRows must be greater than zero")
	}
	if query.DateGbn == "" {
		if query.DateFrom != "" || query.DateTo != "" {
			return application.OldSpcmSearchResult{}, errors.New("dateGbn is required when dateFrom or dateTo is provided")
		}
		return s.port.OldSpcmSearch(ctx, query)
	}
	if query.DateGbn != "1" && query.DateGbn != "2" {
		return application.OldSpcmSearchResult{}, errors.New("dateGbn must be 1 or 2")
	}
	dateFrom, err := parseOldSpcmSearchDate("dateFrom", query.DateFrom)
	if err != nil {
		return application.OldSpcmSearchResult{}, err
	}
	dateTo, err := parseOldSpcmSearchDate("dateTo", query.DateTo)
	if err != nil {
		return application.OldSpcmSearchResult{}, err
	}
	if dateFrom.After(dateTo) {
		return application.OldSpcmSearchResult{}, errors.New("dateFrom must not be after dateTo")
	}

	return s.port.OldSpcmSearch(ctx, query)
}

func parseOldSpcmSearchDate(name, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required when dateGbn is provided", name)
	}
	date, err := time.Parse("20060102", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be a valid date in yyyyMMdd format", name)
	}
	return date, nil
}
