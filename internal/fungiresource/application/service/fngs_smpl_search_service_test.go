package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

type fngsSmplSearchPortStub struct {
	called bool
}

func (s *fngsSmplSearchPortStub) FngsSmplSearch(context.Context, application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error) {
	s.called = true
	return application.FngsSmplSearchResult{}, nil
}

func TestFngsSmplSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.FngsSmplSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.FngsSmplSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.FngsSmplSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.FngsSmplSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &fngsSmplSearchPortStub{}
			service := NewFngsSmplSearchService(port)

			_, err := service.FngsSmplSearch(context.Background(), test.query)
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
