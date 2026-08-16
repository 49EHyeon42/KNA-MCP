package service

import (
	"context"
	"testing"

	"kna-mcp/internal/application"
)

type plantSpecimenSearchPortStub struct {
	called bool
}

func (s *plantSpecimenSearchPortStub) PlantSpecimenSearch(context.Context, application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error) {
	s.called = true
	return application.PlantSpecimenSearchResult{}, nil
}

func TestPlantSpecimenSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSpecimenSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSpecimenSearchQuery{NumberOfRows: 1}, wantError: "pageNumber must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSpecimenSearchQuery{PageNumber: 1}, wantError: "numberOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantSpecimenSearchQuery{PageNumber: 1, NumberOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSpecimenSearchPortStub{}
			service := NewPlantSpecimenSearchService(port)

			_, err := service.PlantSpecimenSearch(context.Background(), test.query)
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
