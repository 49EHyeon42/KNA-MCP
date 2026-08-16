package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSmplSearchPortStub struct {
	called bool
}

func (s *plantSmplSearchPortStub) PlantSmplSearch(context.Context, application.PlantSmplSearchQuery) (application.PlantSmplSearchResult, error) {
	s.called = true
	return application.PlantSmplSearchResult{}, nil
}

func TestPlantSmplSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSmplSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSmplSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSmplSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantSmplSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSmplSearchPortStub{}
			service := NewPlantSmplSearchService(port)

			_, err := service.PlantSmplSearch(context.Background(), test.query)
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
