package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
)

type scnmInfoPortStub struct {
	called bool
}

func (s *scnmInfoPortStub) ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	s.called = true
	return application.ScnmInfoResult{}, nil
}

func TestScnmInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.ScnmInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "empty reqLchnScnmId", wantError: "reqLchnScnmId is required"},
		{name: "blank reqLchnScnmId", query: application.ScnmInfoQuery{ReqLchnScnmID: " \t"}, wantError: "reqLchnScnmId is required"},
		{name: "valid reqLchnScnmId", query: application.ScnmInfoQuery{ReqLchnScnmID: "test-lichen-scientific-name-id"}, wantCall: true},
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
		})
	}
}
