package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
)

type plantMstnsListPortStub struct {
	called bool
}

func (s *plantMstnsListPortStub) PlantMstnsList(context.Context, application.PlantMstnsListQuery) (application.PlantMstnsListResult, error) {
	s.called = true
	return application.PlantMstnsListResult{}, nil
}

func TestPlantMstnsListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantMstnsListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantMstnsListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantMstnsListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantMstnsListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantMstnsListPortStub{}
			service := NewPlantMstnsListService(port)

			_, err := service.PlantMstnsList(context.Background(), test.query)
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
