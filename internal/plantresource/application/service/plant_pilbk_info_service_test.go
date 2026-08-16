package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantPilbkInfoPortStub struct {
	called bool
}

func (s *plantPilbkInfoPortStub) PlantPilbkInfo(context.Context, application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error) {
	s.called = true
	return application.PlantPilbkInfoResult{}, nil
}

func TestPlantPilbkInfoService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantPilbkInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing pictorial book number", wantError: "requestPlantPictorialBookNumber is required"},
		{name: "blank pictorial book number", query: application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "  "}, wantError: "requestPlantPictorialBookNumber is required"},
		{name: "valid pictorial book number", query: application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "test-book-number"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantPilbkInfoPortStub{}
			service := NewPlantPilbkInfoService(port)

			_, err := service.PlantPilbkInfo(context.Background(), test.query)
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
