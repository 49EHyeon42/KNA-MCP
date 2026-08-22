package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

type entogSpcmInfoPortStub struct {
	called bool
}

func (s *entogSpcmInfoPortStub) EntogSpcmInfo(context.Context, application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error) {
	s.called = true
	return application.EntogSpcmInfoResult{}, nil
}

func TestEntogSpcmInfoService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.EntogSpcmInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing q1", wantError: "q1 is required"},
		{name: "blank q1", query: application.EntogSpcmInfoQuery{Q1: " "}, wantError: "q1 is required"},
		{name: "provided q1", query: application.EntogSpcmInfoQuery{Q1: "test-specimen-number"}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &entogSpcmInfoPortStub{}
			service := NewEntogSpcmInfoService(port)

			_, err := service.EntogSpcmInfo(context.Background(), test.query)
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
