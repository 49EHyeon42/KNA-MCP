package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

type fngsSmplUnitListPortStub struct {
	called bool
}

func (s *fngsSmplUnitListPortStub) FngsSmplUnitList(context.Context, application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error) {
	s.called = true
	return application.FngsSmplUnitListResult{}, nil
}

func TestFngsSmplUnitListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.FngsSmplUnitListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.FngsSmplUnitListQuery{NumOfRows: 1, ReqFngsID: "test-fungi-id"}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.FngsSmplUnitListQuery{PageNo: 1, ReqFngsID: "test-fungi-id"}, wantError: "numOfRows must be greater than zero"},
		{name: "blank fungi ID", query: application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "  "}, wantError: "reqFngsId is required"},
		{name: "valid query", query: application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "test-fungi-id"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &fngsSmplUnitListPortStub{}
			service := NewFngsSmplUnitListService(port)

			_, err := service.FngsSmplUnitList(context.Background(), test.query)
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
