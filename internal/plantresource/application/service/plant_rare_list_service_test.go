package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantRareListPortStub struct {
	called bool
}

func (s *plantRareListPortStub) PlantRareList(context.Context, application.PlantRareListQuery) (application.PlantRareListResult, error) {
	s.called = true
	return application.PlantRareListResult{}, nil
}

func TestPlantRareListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantRareListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantRareListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantRareListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantRareListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantRareListPortStub{}
			service := NewPlantRareListService(port)

			_, err := service.PlantRareList(context.Background(), test.query)
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
