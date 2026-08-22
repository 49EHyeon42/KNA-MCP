package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

type entogSpcmSearchPortStub struct {
	called bool
}

func (s *entogSpcmSearchPortStub) EntogSpcmSearch(context.Context, application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error) {
	s.called = true
	return application.EntogSpcmSearchResult{}, nil
}

func TestEntogSpcmSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.EntogSpcmSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid st", query: application.EntogSpcmSearchQuery{St: "5"}, wantError: "st must be one of 1, 2, 3, or 4"},
		{name: "missing sw", query: application.EntogSpcmSearchQuery{St: "1"}, wantError: "sw is required"},
		{name: "blank sw", query: application.EntogSpcmSearchQuery{St: "1", Sw: " "}, wantError: "sw is required"},
		{name: "invalid page number", query: application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "partial Korean name", query: application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}, wantCall: true},
		{name: "partial scientific name", query: application.EntogSpcmSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}, wantCall: true},
		{name: "exact Korean name", query: application.EntogSpcmSearchQuery{St: "3", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}, wantCall: true},
		{name: "exact scientific name", query: application.EntogSpcmSearchQuery{St: "4", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}, wantCall: true},
		{name: "documented dates", query: application.EntogSpcmSearchQuery{St: "2", Sw: "test-search-word", DateGbn: "1", DateFrom: "20200101", DateTo: "20201231", PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &entogSpcmSearchPortStub{}
			service := NewEntogSpcmSearchService(port)

			_, err := service.EntogSpcmSearch(context.Background(), test.query)
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
