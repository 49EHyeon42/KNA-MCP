package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
)

type oldSpcmSearchPortStub struct {
	called bool
}

func (s *oldSpcmSearchPortStub) OldSpcmSearch(context.Context, application.OldSpcmSearchQuery) (application.OldSpcmSearchResult, error) {
	s.called = true
	return application.OldSpcmSearchResult{}, nil
}

func TestOldSpcmSearchService(t *testing.T) {
	valid := application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}
	tests := []struct {
		name      string
		query     application.OldSpcmSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid st", query: application.OldSpcmSearchQuery{St: "3"}, wantError: "st must be 1 or 2"},
		{name: "missing sw", query: application.OldSpcmSearchQuery{St: "1"}, wantError: "sw is required"},
		{name: "invalid page number", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "dates without dateGbn", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateFrom: "20240101"}, wantError: "dateGbn is required when dateFrom or dateTo is provided"},
		{name: "invalid dateGbn", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "3"}, wantError: "dateGbn must be 1 or 2"},
		{name: "missing dateFrom", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateTo: "20241231"}, wantError: "dateFrom is required when dateGbn is provided"},
		{name: "invalid dateFrom", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241301", DateTo: "20241231"}, wantError: "dateFrom must be a valid date in yyyyMMdd format"},
		{name: "missing dateTo", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101"}, wantError: "dateTo is required when dateGbn is provided"},
		{name: "invalid dateTo", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101", DateTo: "20240230"}, wantError: "dateTo must be a valid date in yyyyMMdd format"},
		{name: "reversed dates", query: application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241231", DateTo: "20240101"}, wantError: "dateFrom must not be after dateTo"},
		{name: "valid conditions", query: valid, wantCall: true},
		{name: "valid dates", query: application.OldSpcmSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "2", DateFrom: "20240101", DateTo: "20241231"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &oldSpcmSearchPortStub{}
			service := NewOldSpcmSearchService(port)

			_, err := service.OldSpcmSearch(context.Background(), test.query)
			if test.wantError == "" && err != nil {
				t.Fatal(err)
			}
			if test.wantError != "" && (err == nil || err.Error() != test.wantError) {
				t.Errorf("error = %v, want %q", err, test.wantError)
			}
			if port.called != test.wantCall {
				t.Errorf("port called = %t, want %t", port.called, test.wantCall)
			}
		})
	}
}
