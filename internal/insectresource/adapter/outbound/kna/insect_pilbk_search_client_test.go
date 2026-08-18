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

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

func TestInsectPilbkSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/InsectService/insectPilbkSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/InsectService/insectPilbkSearch")
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "2",
			"numOfRows":    "1",
			"reqSearchWrd": "test-search-word",
			"dateFrom":     "test-date-from",
			"dateTo":       "test-date-to",
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
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <genusKorNm>genus Korean name</genusKorNm><genusNm>genus name</genusNm>
      <insctGnrlNm>insect general name</insctGnrlNm><insctPilbkNo>insect pictorial book number</insctPilbkNo>
      <insctSpecsScnm>insect species scientific name</insctSpecsScnm><lastUpdtDtm> </lastUpdtDtm>
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

	got, err := client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{
		PageNo:       2,
		NumOfRows:    1,
		ReqSearchWrd: "test-search-word",
		DateFrom:     "test-date-from",
		DateTo:       "test-date-to",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.InsectPilbkSearchResult{
		Items: []application.InsectPilbkSearchItem{{
			FamilyKorNm:    "family Korean name",
			FamilyNm:       "family name",
			GenusKorNm:     "genus Korean name",
			GenusNm:        "genus name",
			InsctGnrlNm:    "insect general name",
			InsctPilbkNo:   "insect pictorial book number",
			InsctSpecsScnm: "insect species scientific name",
			LastUpdtDtm:    " ",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestInsectPilbkSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, key := range []string{"reqSearchWrd", "dateFrom", "dateTo"} {
			if _, exists := request.URL.Query()[key]; exists {
				t.Errorf("%s query was sent for an empty value", key)
			}
		}
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestInsectPilbkSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *InsectPilbkSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *InsectPilbkSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestInsectPilbkSearchReturnsObservedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>99</resultCode><resultMsg>ORA-00908: missing NULL keyword</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1, DateFrom: "20240101"})
	var apiError *InsectPilbkSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *InsectPilbkSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusOK || apiError.Code != "99" || apiError.Message != "ORA-00908: missing NULL keyword" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestInsectPilbkSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
	var apiError *InsectPilbkSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *InsectPilbkSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestInsectPilbkSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "insectPilbkSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "insectPilbkSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "insectPilbkSearch: response missing resultCode"},
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

			_, err = client.InsectPilbkSearch(context.Background(), application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestInsectPilbkSearchLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name           string
		reqSearchWrd   string
		insctGnrlNm    string
		insctSpecsScnm string
	}{
		{name: "without search word"},
		{name: "exact Korean name", reqSearchWrd: "가락지나비", insctGnrlNm: "가락지나비"},
		{name: "partial Korean name", reqSearchWrd: "가락지", insctGnrlNm: "가락지나비"},
		{name: "uppercase scientific name", reqSearchWrd: "Aphantopus", insctSpecsScnm: "aphantopus"},
		{name: "lowercase scientific name", reqSearchWrd: "aphantopus", insctSpecsScnm: "aphantopus"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.InsectPilbkSearch(ctx, application.InsectPilbkSearchQuery{
				PageNo:       1,
				NumOfRows:    10,
				ReqSearchWrd: test.reqSearchWrd,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("insectPilbkSearch returned no items")
			}
			if test.insctGnrlNm != "" || test.insctSpecsScnm != "" {
				matched := false
				for _, item := range result.Items {
					if test.insctGnrlNm != "" && item.InsctGnrlNm == test.insctGnrlNm {
						matched = true
					}
					if test.insctSpecsScnm != "" && strings.Contains(strings.ToLower(item.InsctSpecsScnm), test.insctSpecsScnm) {
						matched = true
					}
				}
				if !matched {
					t.Errorf("insectPilbkSearch items = %#v, want matching result", result.Items)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.InsectPilbkSearch(ctx, application.InsectPilbkSearchQuery{PageNo: 1, NumOfRows: 2, ReqSearchWrd: "나비"})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.InsectPilbkSearch(ctx, application.InsectPilbkSearchQuery{PageNo: 2, NumOfRows: 2, ReqSearchWrd: "나비"})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 2 || len(second.Items) != 2 || first.TotalCount != second.TotalCount || first.TotalCount < 4 || first.Items[0].InsctPilbkNo == second.Items[0].InsctPilbkNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.InsectPilbkSearch(ctx, application.InsectPilbkSearchQuery{
			PageNo:       1,
			NumOfRows:    1,
			ReqSearchWrd: "kna-mcp-no-result-20260817",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
