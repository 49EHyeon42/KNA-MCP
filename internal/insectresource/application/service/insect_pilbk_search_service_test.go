package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

type insectPilbkSearchPortStub struct {
	called bool
}

func (s *insectPilbkSearchPortStub) InsectPilbkSearch(context.Context, application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error) {
	s.called = true
	return application.InsectPilbkSearchResult{}, nil
}

func TestInsectPilbkSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.InsectPilbkSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.InsectPilbkSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.InsectPilbkSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &insectPilbkSearchPortStub{}
			service := NewInsectPilbkSearchService(port)

			_, err := service.InsectPilbkSearch(context.Background(), test.query)
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
