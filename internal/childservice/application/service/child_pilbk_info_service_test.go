package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

type childPilbkInfoPortStub struct {
	called bool
}

func (s *childPilbkInfoPortStub) ChildPilbkInfo(context.Context, application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error) {
	s.called = true
	return application.ChildPilbkInfoResult{}, nil
}

func TestChildPilbkInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.ChildPilbkInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing child pictorial book number", wantError: "reqChildLvbngPilbkNo is required"},
		{name: "blank child pictorial book number", query: application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: " \t"}, wantError: "reqChildLvbngPilbkNo is required"},
		{name: "valid child pictorial book number", query: application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"}, wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &childPilbkInfoPortStub{}
			service := NewChildPilbkInfoService(port)

			_, err := service.ChildPilbkInfo(context.Background(), test.query)
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
