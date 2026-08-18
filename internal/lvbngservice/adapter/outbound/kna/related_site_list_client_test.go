package kna

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
)

func TestRelatedSiteList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/LvbngService2/relatedSiteList" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/LvbngService2/relatedSiteList")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"pageNo":     "2",
			"numOfRows":  "1",
		}
		if len(query) != len(wantQuery) {
			t.Errorf("query key count = %d, want %d", len(query), len(wantQuery))
		}
		for key, want := range wantQuery {
			if got := query.Get(key); got != want {
				t.Errorf("query %s = %q, want %q", key, got, want)
			}
		}

		response.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(response, `<?xml version="1.0" encoding="UTF-8"?>
<response>
  <header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header>
  <body>
    <items><item>
      <lvbngTpcdNm>living thing type code name</lvbngTpcdNm>
      <siteCtgryNm> </siteCtgryNm>
      <siteNm>site name</siteNm>
      <siteUrl>http://example.com</siteUrl>
    </item></items>
    <numOfRows>1</numOfRows><pageNo>2</pageNo><totalCount>7</totalCount>
  </body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.RelatedSiteList(context.Background(), application.RelatedSiteListQuery{PageNo: 2, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}

	want := application.RelatedSiteListResult{
		Items: []application.RelatedSiteListItem{{
			LvbngTpcdNm: "living thing type code name",
			SiteCtgryNm: " ",
			SiteNm:      "site name",
			SiteURL:     "http://example.com",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestRelatedSiteListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>10</numOfRows><pageNo>999</pageNo><totalCount>7</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.RelatedSiteList(context.Background(), application.RelatedSiteListQuery{PageNo: 999, NumOfRows: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 7 {
		t.Errorf("result = %#v, want empty page with totalCount 7", result)
	}
}

func TestRelatedSiteListReturnsDocumentedAPIErrors(t *testing.T) {
	tests := []struct {
		code    string
		message string
	}{
		{code: "02", message: "DB_ERROR"},
		{code: "03", message: "NODATA_ERROR"},
		{code: "05", message: "SERVICETIME_OUT"},
		{code: "10", message: "INVALID_REQUEST_PARAMETER_ERROR"},
		{code: "11", message: "NO_MANDATORY_REQUEST_PARAMETERS_ERROR"},
		{code: "21", message: "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR"},
		{code: "33", message: "UNSIGNED_CALL_ERROR"},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(response, `<response><header><resultCode>%s</resultCode></header></response>`, test.code)
			}))
			defer server.Close()

			client, err := NewClient("test-key")
			if err != nil {
				t.Fatal(err)
			}
			client.baseURL = server.URL

			_, err = client.RelatedSiteList(context.Background(), application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 1})
			var apiError *RelatedSiteListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *RelatedSiteListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestRelatedSiteListReturnsGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(response, `<OpenAPI_ServiceResponse><cmmMsgHeader><errMsg>SERVICE_KEY_IS_NOT_REGISTERED_ERROR</errMsg><returnAuthMsg>등록되지 않은 서비스키</returnAuthMsg><returnReasonCode>30</returnReasonCode></cmmMsgHeader></OpenAPI_ServiceResponse>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.RelatedSiteList(context.Background(), application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 1})
	var apiError *RelatedSiteListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *RelatedSiteListError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestRelatedSiteListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "relatedSiteList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "relatedSiteList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "relatedSiteList: response missing resultCode"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				response.WriteHeader(test.statusCode)
				_, _ = io.WriteString(response, test.body)
			}))
			defer server.Close()

			client, err := NewClient("test-key")
			if err != nil {
				t.Fatal(err)
			}
			client.baseURL = server.URL

			_, err = client.RelatedSiteList(context.Background(), application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestRelatedSiteListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("distinct pages", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.RelatedSiteList(ctx, application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.RelatedSiteList(ctx, application.RelatedSiteListQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || reflect.DeepEqual(first.Items[0], second.Items[0]) {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("preserves documented fields", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.RelatedSiteList(ctx, application.RelatedSiteListQuery{PageNo: 1, NumOfRows: 100})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != result.TotalCount || result.TotalCount == 0 {
			t.Fatalf("result = %#v, want all related sites", result)
		}
		spaceCategory := false
		for _, item := range result.Items {
			if item.LvbngTpcdNm == "" || item.SiteNm == "" || item.SiteURL == "" {
				t.Errorf("item = %#v, want documented fields", item)
			}
			if !strings.HasPrefix(item.SiteURL, "http://") && !strings.HasPrefix(item.SiteURL, "https://") {
				t.Errorf("siteUrl = %q, want original HTTP or HTTPS URL", item.SiteURL)
			}
			if item.SiteCtgryNm == " " {
				spaceCategory = true
			}
		}
		if !spaceCategory {
			t.Error("siteCtgryNm did not preserve the observed single-space value")
		}
	})

	t.Run("empty page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.RelatedSiteList(ctx, application.RelatedSiteListQuery{PageNo: 999, NumOfRows: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount == 0 {
			t.Errorf("result = %#v, want empty page with nonzero totalCount", result)
		}
	})
}
