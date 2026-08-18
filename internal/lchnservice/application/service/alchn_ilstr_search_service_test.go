package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

type alchnIlstrSearchPortStub struct {
	called bool
}

func (s *alchnIlstrSearchPortStub) AlchnIlstrSearch(context.Context, application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error) {
	s.called = true
	return application.AlchnIlstrSearchResult{}, nil
}

func TestAlchnIlstrSearchService(t *testing.T) {
	valid := application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1}
	tests := []struct {
		name      string
		query     application.AlchnIlstrSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid st", query: application.AlchnIlstrSearchQuery{St: "5"}, wantError: "st must be one of 1, 2, 3, or 4"},
		{name: "missing sw", query: application.AlchnIlstrSearchQuery{St: "2"}, wantError: "sw is required"},
		{name: "invalid page number", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "dates without dateGbn", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateFrom: "20240101"}, wantError: "dateGbn is required when dateFrom or dateTo is provided"},
		{name: "invalid dateGbn", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "3"}, wantError: "dateGbn must be 1 or 2"},
		{name: "missing dateFrom", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateTo: "20241231"}, wantError: "dateFrom is required when dateGbn is provided"},
		{name: "invalid dateFrom", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241301", DateTo: "20241231"}, wantError: "dateFrom must be a valid date in yyyyMMdd format"},
		{name: "missing dateTo", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101"}, wantError: "dateTo is required when dateGbn is provided"},
		{name: "invalid dateTo", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240101", DateTo: "20240230"}, wantError: "dateTo must be a valid date in yyyyMMdd format"},
		{name: "reversed dates", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241231", DateTo: "20240101"}, wantError: "dateFrom must not be after dateTo"},
		{name: "valid conditions", query: valid, wantCall: true},
		{name: "valid dates", query: application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 1, DateGbn: "2", DateFrom: "20240101", DateTo: "20241231"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &alchnIlstrSearchPortStub{}
			service := NewAlchnIlstrSearchService(port)

			_, err := service.AlchnIlstrSearch(context.Background(), test.query)
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
