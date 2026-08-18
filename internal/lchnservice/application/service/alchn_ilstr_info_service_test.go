package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

type alchnIlstrInfoPortStub struct {
	called bool
}

func (s *alchnIlstrInfoPortStub) AlchnIlstrInfo(context.Context, application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error) {
	s.called = true
	return application.AlchnIlstrInfoResult{}, nil
}

func TestAlchnIlstrInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		q1        string
		wantError string
		wantCall  bool
	}{
		{name: "missing q1", wantError: "q1 is required"},
		{name: "blank q1", q1: " ", wantError: "q1 is required"},
		{name: "valid q1", q1: "test-lichen-pictorial-book-number", wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &alchnIlstrInfoPortStub{}
			service := NewAlchnIlstrInfoService(port)

			_, err := service.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: test.q1})
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
