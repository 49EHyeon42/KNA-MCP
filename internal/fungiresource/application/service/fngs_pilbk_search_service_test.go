package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

type fngsPilbkSearchPortStub struct {
	called bool
}

func (s *fngsPilbkSearchPortStub) FngsPilbkSearch(context.Context, application.FngsPilbkSearchQuery) (application.FngsPilbkSearchResult, error) {
	s.called = true
	return application.FngsPilbkSearchResult{}, nil
}

func TestFngsPilbkSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.FngsPilbkSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.FngsPilbkSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.FngsPilbkSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.FngsPilbkSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &fngsPilbkSearchPortStub{}
			service := NewFngsPilbkSearchService(port)

			_, err := service.FngsPilbkSearch(context.Background(), test.query)
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
