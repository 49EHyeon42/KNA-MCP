package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSeedSearchPortStub struct {
	called bool
}

func (s *plantSeedSearchPortStub) PlantSeedSearch(context.Context, application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error) {
	s.called = true
	return application.PlantSeedSearchResult{}, nil
}

func TestPlantSeedSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSeedSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSeedSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSeedSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSeedSearchPortStub{}
			service := NewPlantSeedSearchService(port)

			_, err := service.PlantSeedSearch(context.Background(), test.query)
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
