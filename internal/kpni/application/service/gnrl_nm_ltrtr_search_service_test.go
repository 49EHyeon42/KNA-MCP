package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

type gnrlNmLtrtrSearchPortStub struct {
	called bool
	query  application.GnrlNmLtrtrSearchQuery
}

func (s *gnrlNmLtrtrSearchPortStub) GnrlNmLtrtrSearch(_ context.Context, query application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error) {
	s.called = true
	s.query = query
	return application.GnrlNmLtrtrSearchResult{}, nil
}

func TestGnrlNmLtrtrSearchService(t *testing.T) {
	for _, test := range []struct {
		name      string
		query     application.GnrlNmLtrtrSearchQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.GnrlNmLtrtrSearchQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.GnrlNmLtrtrSearchQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid conditions", query: application.GnrlNmLtrtrSearchQuery{PageNo: 1, NumOfRows: 10, ReqPlantGnrlNm: "소나무"}, wantCall: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			port := &gnrlNmLtrtrSearchPortStub{}
			service := NewGnrlNmLtrtrSearchService(port)

			_, err := service.GnrlNmLtrtrSearch(context.Background(), test.query)
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
