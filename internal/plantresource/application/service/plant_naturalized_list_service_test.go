package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantNaturalizedListPortStub struct {
	called bool
}

func (s *plantNaturalizedListPortStub) PlantNaturalizedList(context.Context, application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error) {
	s.called = true
	return application.PlantNaturalizedListResult{}, nil
}

func TestPlantNaturalizedListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantNaturalizedListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantNaturalizedListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantNaturalizedListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantNaturalizedListPortStub{}
			service := NewPlantNaturalizedListService(port)

			_, err := service.PlantNaturalizedList(context.Background(), test.query)
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
