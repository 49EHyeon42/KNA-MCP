package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

type insectSmplSearchPortStub struct {
	called bool
}

func (s *insectSmplSearchPortStub) InsectSmplSearch(context.Context, application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error) {
	s.called = true
	return application.InsectSmplSearchResult{}, nil
}

func TestInsectSmplSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.InsectSmplSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.InsectSmplSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.InsectSmplSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.InsectSmplSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &insectSmplSearchPortStub{}
			service := NewInsectSmplSearchService(port)

			_, err := service.InsectSmplSearch(context.Background(), test.query)
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
