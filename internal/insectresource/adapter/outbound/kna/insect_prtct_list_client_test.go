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

func TestInsectPrtctList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/InsectService/insectPrtctList" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/InsectService/insectPrtctList")
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
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <insctGnrlNm>insect general name</insctGnrlNm><insctPcmtt>endangered classification</insctPcmtt>
      <insctPilbkNo>insect pictorial book number</insctPilbkNo>
      <insctSpecsScnm>insect species scientific name</insctSpecsScnm>
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

	got, err := client.InsectPrtctList(context.Background(), application.InsectPrtctListQuery{PageNo: 2, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}

	want := application.InsectPrtctListResult{
		Items: []application.InsectPrtctListItem{{
			FamilyKorNm:    "family Korean name",
			FamilyNm:       "family name",
			InsctGnrlNm:    "insect general name",
			InsctPcmtt:     "endangered classification",
			InsctPilbkNo:   "insect pictorial book number",
			InsctSpecsScnm: "insect species scientific name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestInsectPrtctListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1000000</pageNo><totalCount>7</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.InsectPrtctList(context.Background(), application.InsectPrtctListQuery{PageNo: 1000000, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 7 {
		t.Errorf("result = %#v, want empty items and totalCount 7", result)
	}
}

func TestInsectPrtctListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.InsectPrtctList(context.Background(), application.InsectPrtctListQuery{PageNo: 1, NumOfRows: 1})
			var apiError *InsectPrtctListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *InsectPrtctListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestInsectPrtctListReturnsGatewayError(t *testing.T) {
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

	_, err = client.InsectPrtctList(context.Background(), application.InsectPrtctListQuery{PageNo: 1, NumOfRows: 1})
	var apiError *InsectPrtctListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *InsectPrtctListError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestInsectPrtctListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "insectPrtctList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "insectPrtctList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "insectPrtctList: response missing resultCode"},
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

			_, err = client.InsectPrtctList(context.Background(), application.InsectPrtctListQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestInsectPrtctListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.InsectPrtctList(ctx, application.InsectPrtctListQuery{PageNo: 1, NumOfRows: 5})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.InsectPrtctList(ctx, application.InsectPrtctListQuery{PageNo: 2, NumOfRows: 5})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 5 || len(second.Items) != 5 || first.TotalCount != second.TotalCount || first.TotalCount < 10 || first.Items[0].InsctPilbkNo == second.Items[0].InsctPilbkNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
		for _, item := range append(first.Items, second.Items...) {
			if item.InsctPilbkNo == "" || item.InsctGnrlNm == "" || item.InsctPcmtt == "" {
				t.Errorf("item = %#v, want pictorial book number, general name, and endangered classification", item)
			}
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.InsectPrtctList(ctx, application.InsectPrtctListQuery{PageNo: 1000000, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 {
			t.Errorf("result = %#v, want empty items", result)
		}
	})
}
