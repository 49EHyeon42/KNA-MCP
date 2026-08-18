package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

type alchnSpcmSearchPortStub struct {
	called bool
}

func (s *alchnSpcmSearchPortStub) AlchnSpcmSearch(context.Context, application.AlchnSpcmSearchQuery) (application.AlchnSpcmSearchResult, error) {
	s.called = true
	return application.AlchnSpcmSearchResult{}, nil
}

func TestAlchnSpcmSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.AlchnSpcmSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid st", query: application.AlchnSpcmSearchQuery{St: "5"}, wantError: "st must be one of 1, 2, 3, or 4"},
		{name: "missing sw", query: application.AlchnSpcmSearchQuery{St: "2"}, wantError: "sw is required"},
		{name: "invalid page number", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "dates without dateGbn", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateFrom: "20240101"}, wantError: "dateGbn is required when dateFrom or dateTo is provided"},
		{name: "invalid dateGbn", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "3"}, wantError: "dateGbn must be 1 or 2"},
		{name: "missing dateFrom", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateTo: "20241231"}, wantError: "dateFrom is required when dateGbn is provided"},
		{name: "invalid dateFrom", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241301", DateTo: "20241231"}, wantError: "dateFrom must be a valid date in yyyyMMdd format"},
		{name: "missing dateTo", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101"}, wantError: "dateTo is required when dateGbn is provided"},
		{name: "invalid dateTo", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101", DateTo: "20240230"}, wantError: "dateTo must be a valid date in yyyyMMdd format"},
		{name: "reversed dates", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241231", DateTo: "20240101"}, wantError: "dateFrom must not be after dateTo"},
		{name: "valid conditions", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1}, wantCall: true},
		{name: "valid dates", query: application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "2", DateFrom: "20240101", DateTo: "20241231"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &alchnSpcmSearchPortStub{}
			service := NewAlchnSpcmSearchService(port)

			_, err := service.AlchnSpcmSearch(context.Background(), test.query)
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
