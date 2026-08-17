package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

type insectPrtctListPortStub struct {
	called bool
}

func (s *insectPrtctListPortStub) InsectPrtctList(context.Context, application.InsectPrtctListQuery) (application.InsectPrtctListResult, error) {
	s.called = true
	return application.InsectPrtctListResult{}, nil
}

func TestInsectPrtctListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.InsectPrtctListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.InsectPrtctListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.InsectPrtctListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid query", query: application.InsectPrtctListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &insectPrtctListPortStub{}
			service := NewInsectPrtctListService(port)

			_, err := service.InsectPrtctList(context.Background(), test.query)
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
