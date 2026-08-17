package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

type insectPilbkInfoPortStub struct {
	called bool
	query  application.InsectPilbkInfoQuery
}

func (s *insectPilbkInfoPortStub) InsectPilbkInfo(_ context.Context, query application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error) {
	s.called = true
	s.query = query
	return application.InsectPilbkInfoResult{}, nil
}

func TestInsectPilbkInfoService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.InsectPilbkInfoQuery
		wantError string
		wantCall  bool
	}{
		{name: "missing insect pictorial book number", wantError: "reqInsctPilbkNo is required"},
		{name: "blank insect pictorial book number", query: application.InsectPilbkInfoQuery{ReqInsctPilbkNo: " \t"}, wantError: "reqInsctPilbkNo is required"},
		{name: "valid condition", query: application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"}, wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &insectPilbkInfoPortStub{}
			service := NewInsectPilbkInfoService(port)

			_, err := service.InsectPilbkInfo(context.Background(), test.query)
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
