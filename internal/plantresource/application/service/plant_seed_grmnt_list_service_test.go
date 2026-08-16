package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSeedGrmntListPortStub struct {
	called bool
}

func (s *plantSeedGrmntListPortStub) PlantSeedGrmntList(context.Context, application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error) {
	s.called = true
	return application.PlantSeedGrmntListResult{}, nil
}

func TestPlantSeedGrmntListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantSeedGrmntListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantSeedGrmntListQuery{NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantSeedGrmntListQuery{PageNo: 1, ReqSeedSpecsID: "test-seed-species-id"}, wantError: "numOfRows must be greater than zero"},
		{name: "blank seed species ID", query: application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "  "}, wantError: "reqSeedSpecsId is required"},
		{name: "valid query", query: application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantSeedGrmntListPortStub{}
			service := NewPlantSeedGrmntListService(port)

			_, err := service.PlantSeedGrmntList(context.Background(), test.query)
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
