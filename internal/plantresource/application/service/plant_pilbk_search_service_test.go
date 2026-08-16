package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantPilbkSearchPortStub struct {
	called bool
}

func (s *plantPilbkSearchPortStub) PlantPilbkSearch(context.Context, application.PlantPilbkSearchQuery) (application.PlantPilbkSearchResult, error) {
	s.called = true
	return application.PlantPilbkSearchResult{}, nil
}

func TestPlantPilbkSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantPilbkSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantPilbkSearchQuery{NumOfRows: 1}, wantError: "pageNumber must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantPilbkSearchQuery{PageNo: 1}, wantError: "numberOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantPilbkSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantPilbkSearchPortStub{}
			service := NewPlantPilbkSearchService(port)

			_, err := service.PlantPilbkSearch(context.Background(), test.query)
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
