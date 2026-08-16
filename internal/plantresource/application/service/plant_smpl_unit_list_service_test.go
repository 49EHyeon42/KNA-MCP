package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSmplUnitListPortStub struct {
	called bool
}

func (s *plantSmplUnitListPortStub) PlantSmplUnitList(context.Context, application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error) {
	s.called = true
	return application.PlantSmplUnitListResult{}, nil
}

func TestPlantSmplUnitListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSmplUnitListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSmplUnitListQuery{NumOfRows: 1, ReqPlantSpecsID: "test-plant-species-id"}, wantError: "pageNumber must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSmplUnitListQuery{PageNo: 1, ReqPlantSpecsID: "test-plant-species-id"}, wantError: "numberOfRows must be greater than zero"},
		{name: "blank plant species ID", query: application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqPlantSpecsID: "  "}, wantError: "requestPlantSpeciesId is required"},
		{name: "valid query", query: application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqPlantSpecsID: "test-plant-species-id"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSmplUnitListPortStub{}
			service := NewPlantSmplUnitListService(port)

			_, err := service.PlantSmplUnitList(context.Background(), test.query)
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
