package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

type fngsPilbkInfoPortStub struct {
	called bool
}

func (s *fngsPilbkInfoPortStub) FngsPilbkInfo(context.Context, application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error) {
	s.called = true
	return application.FngsPilbkInfoResult{}, nil
}

func TestFngsPilbkInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.FngsPilbkInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing fungi pictorial book number", wantError: "reqFngsPilbkNo is required"},
		{name: "blank fungi pictorial book number", query: application.FngsPilbkInfoQuery{ReqFngsPilbkNo: " \t"}, wantError: "reqFngsPilbkNo is required"},
		{name: "valid fungi pictorial book number", query: application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"}, wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &fngsPilbkInfoPortStub{}
			service := NewFngsPilbkInfoService(port)

			_, err := service.FngsPilbkInfo(context.Background(), test.query)
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
