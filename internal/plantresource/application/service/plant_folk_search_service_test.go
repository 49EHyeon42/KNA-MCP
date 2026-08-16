package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantFolkSearchPortStub struct {
	called bool
}

func (s *plantFolkSearchPortStub) PlantFolkSearch(context.Context, application.PlantFolkSearchQuery) (application.PlantFolkSearchResult, error) {
	s.called = true
	return application.PlantFolkSearchResult{}, nil
}

func TestPlantFolkSearchService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.PlantFolkSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.PlantFolkSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.PlantFolkSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &plantFolkSearchPortStub{}
			service := NewPlantFolkSearchService(port)

			_, err := service.PlantFolkSearch(context.Background(), test.query)
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
