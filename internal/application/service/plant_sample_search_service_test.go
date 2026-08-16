package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

type plantSampleSearchPortStub struct {
	called bool
}

func (s *plantSampleSearchPortStub) PlantSampleSearch(context.Context, application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error) {
	s.called = true
	return application.PlantSampleSearchResult{}, nil
}

func TestPlantSampleSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSampleSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSampleSearchQuery{NumberOfRows: 1}, wantError: "pageNumber must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSampleSearchQuery{PageNumber: 1}, wantError: "numberOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantSampleSearchQuery{PageNumber: 1, NumberOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSampleSearchPortStub{}
			service := NewPlantSampleSearchService(port)

			_, err := service.PlantSampleSearch(context.Background(), test.query)
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
