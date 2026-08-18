package service

import (
	"context"
	"testing"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
)

type relatedSiteListPortStub struct {
	called bool
}

func (s *relatedSiteListPortStub) RelatedSiteList(context.Context, application.RelatedSiteListQuery) (application.RelatedSiteListResult, error) {
	s.called = true
	return application.RelatedSiteListResult{}, nil
}

func TestRelatedSiteListService(t *testing.T) {
	tests := []struct {
		name      string
		query     application.RelatedSiteListQuery
		wantError string
		wantCall  bool
	}{
		{name: "invalid page number", query: application.RelatedSiteListQuery{NumOfRows: 1}, wantError: "pageNo must be greater than zero"},
		{name: "invalid number of rows", query: application.RelatedSiteListQuery{PageNo: 1}, wantError: "numOfRows must be greater than zero"},
		{name: "valid pagination", query: application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 1}, wantCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := &relatedSiteListPortStub{}
			service := NewRelatedSiteListService(port)

			_, err := service.RelatedSiteList(context.Background(), test.query)
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
