package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

type alchnSpcmInfoPortStub struct {
	called bool
}

func (s *alchnSpcmInfoPortStub) AlchnSpcmInfo(context.Context, application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error) {
	s.called = true
	return application.AlchnSpcmInfoResult{}, nil
}

func TestAlchnSpcmInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		q1        string
		wantError string
		wantCall  bool
	}{
		{name: "missing q1", wantError: "q1 is required"},
		{name: "blank q1", q1: " ", wantError: "q1 is required"},
		{name: "valid q1", q1: "TEST-SAMPLE-001", wantCall: true},
		{name: "unrecognized nonblank q1", q1: "UNRECOGNIZED", wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &alchnSpcmInfoPortStub{}
			service := NewAlchnSpcmInfoService(port)

			_, err := service.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: test.q1})
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
