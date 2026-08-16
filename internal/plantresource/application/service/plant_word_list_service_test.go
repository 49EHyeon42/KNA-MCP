package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantWordListPortStub struct {
	called bool
}

func (s *plantWordListPortStub) PlantWordList(context.Context, application.PlantWordListQuery) (application.PlantWordListResult, error) {
	s.called = true
	return application.PlantWordListResult{}, nil
}

func TestPlantWordListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantWordListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantWordListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantWordListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantWordListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantWordListPortStub{}
			service := NewPlantWordListService(port)

			_, err := service.PlantWordList(context.Background(), test.query)
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
