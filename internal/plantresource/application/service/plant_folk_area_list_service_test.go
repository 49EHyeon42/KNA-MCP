package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantFolkAreaListPortStub struct {
	called bool
}

func (s *plantFolkAreaListPortStub) PlantFolkAreaList(context.Context, application.PlantFolkAreaListQuery) (application.PlantFolkAreaListResult, error) {
	s.called = true
	return application.PlantFolkAreaListResult{}, nil
}

func TestPlantFolkAreaListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantFolkAreaListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantFolkAreaListQuery{NumOfRows: 1, FlpltID: "test-folk-plant-id"}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantFolkAreaListQuery{PageNo: 1, FlpltID: "test-folk-plant-id"}, wantError: "numOfRows must be greater than zero"},
		{name: "blank folk plant ID", query: application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "  "}, wantError: "flpltId is required"},
		{name: "valid query", query: application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "test-folk-plant-id"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantFolkAreaListPortStub{}
			service := NewPlantFolkAreaListService(port)

			_, err := service.PlantFolkAreaList(context.Background(), test.query)
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
