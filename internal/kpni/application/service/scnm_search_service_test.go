package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

type scnmSearchPortStub struct {
	called bool
}

func (s *scnmSearchPortStub) ScnmSearch(context.Context, application.ScnmSearchQuery) (application.ScnmSearchResult, error) {
	s.called = true
	return application.ScnmSearchResult{}, nil
}

func TestScnmSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.ScnmSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.ScnmSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.ScnmSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "invalid dateFrom format", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1, DateFrom: "2024/01/01"}, wantError: "dateFrom must be a valid date in yyyyMMdd format"},
		{name: "invalid dateTo value", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1, DateTo: "20240230"}, wantError: "dateTo must be a valid date in yyyyMMdd format"},
		{name: "valid conditions", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1, DateFrom: "20240101", DateTo: "20241231"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &scnmSearchPortStub{}
			service := NewScnmSearchService(port)

			_, err := service.ScnmSearch(context.Background(), test.query)
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
