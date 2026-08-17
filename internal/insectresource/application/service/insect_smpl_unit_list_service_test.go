package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

type insectSmplUnitListPortStub struct {
	called bool
}

func (s *insectSmplUnitListPortStub) InsectSmplUnitList(context.Context, application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error) {
	s.called = true
	return application.InsectSmplUnitListResult{}, nil
}

func TestInsectSmplUnitListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.InsectSmplUnitListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.InsectSmplUnitListQuery{NumOfRows: 1, ReqInsctSpecsID: "test-insect-species-id"}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.InsectSmplUnitListQuery{PageNo: 1, ReqInsctSpecsID: "test-insect-species-id"}, wantError: "numOfRows must be greater than zero"},
		{name: "blank insect species ID", query: application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "  "}, wantError: "reqInsctSpecsId is required"},
		{name: "valid query", query: application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "test-insect-species-id"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &insectSmplUnitListPortStub{}
			service := NewInsectSmplUnitListService(port)

			_, err := service.InsectSmplUnitList(context.Background(), test.query)
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
