package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

type scnmInfoPortStub struct {
	called bool
	query  application.ScnmInfoQuery
}

func (s *scnmInfoPortStub) ScnmInfo(_ context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	s.called = true
	s.query = query
	return application.ScnmInfoResult{}, nil
}

func TestScnmInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.ScnmInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing scientific name ID", wantError: "reqPlantScnmId is required"},
		{name: "blank scientific name ID", query: application.ScnmInfoQuery{ReqPlantScnmID: " \t"}, wantError: "reqPlantScnmId is required"},
		{name: "valid condition", query: application.ScnmInfoQuery{ReqPlantScnmID: "test-plant-scientific-name-id"}, wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &scnmInfoPortStub{}
			service := NewScnmInfoService(port)

			_, err := service.ScnmInfo(context.Background(), test.query)
			if test.wantError == "" && err != nil {
				t.Fatal(err)
			}
			if test.wantError != "" && (err == nil || err.Error() != test.wantError) {
				t.Errorf("error = %v, want %q", err, test.wantError)
			}
			if port.called != test.wantCall {
				t.Errorf("port called = %t, want %t", port.called, test.wantCall)
			}
			if port.called && port.query != test.query {
				t.Errorf("query = %#v, want %#v", port.query, test.query)
			}
		})
	}
}
