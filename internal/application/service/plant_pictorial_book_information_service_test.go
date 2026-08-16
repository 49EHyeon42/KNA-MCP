package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

type plantPictorialBookInformationPortStub struct {
	called bool
}

func (s *plantPictorialBookInformationPortStub) PlantPictorialBookInformation(context.Context, application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error) {
	s.called = true
	return application.PlantPictorialBookInformationResult{}, nil
}

func TestPlantPictorialBookInformationService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantPictorialBookInformationQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing pictorial book number", wantError: "requestPlantPictorialBookNumber is required"},
		{name: "blank pictorial book number", query: application.PlantPictorialBookInformationQuery{RequestPlantPictorialBookNumber: "  "}, wantError: "requestPlantPictorialBookNumber is required"},
		{name: "valid pictorial book number", query: application.PlantPictorialBookInformationQuery{RequestPlantPictorialBookNumber: "test-book-number"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantPictorialBookInformationPortStub{}
			service := NewPlantPictorialBookInformationService(port)

			_, err := service.PlantPictorialBookInformation(context.Background(), test.query)
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
