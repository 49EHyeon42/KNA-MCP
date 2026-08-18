package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

type childPilbkSearchPortStub struct {
	called bool
}

func (s *childPilbkSearchPortStub) ChildPilbkSearch(context.Context, application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error) {
	s.called = true
	return application.ChildPilbkSearchResult{}, nil
}

func TestChildPilbkSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.ChildPilbkSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.ChildPilbkSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.ChildPilbkSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &childPilbkSearchPortStub{}
			service := NewChildPilbkSearchService(port)

			_, err := service.ChildPilbkSearch(context.Background(), test.query)
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
